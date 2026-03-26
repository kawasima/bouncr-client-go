package bouncr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestServer creates a test server that handles OAuth2 token requests
// and routes API requests to the given handler.
func newTestServer(t *testing.T, handler func(http.ResponseWriter, *http.Request)) (*httptest.Server, *Client) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth2/token" {
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}
		handler(res, req)
	}))
	client, err := NewClientWithOptions("test-client-id", "test-client-secret", ts.URL, false)
	if err != nil {
		t.Fatal(err)
	}
	return ts, client
}

func TestRequest(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	})
	defer ts.Close()

	req, _ := http.NewRequestWithContext(context.Background(), "GET", client.urlFor("/test").String(), nil)
	resp, err := client.Request(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRequestAPIError(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
	})
	defer ts.Close()

	req, _ := http.NewRequestWithContext(context.Background(), "GET", client.urlFor("/missing").String(), nil)
	_, err := client.Request(context.Background(), req)
	if err == nil {
		t.Error("expected error for 404 response")
	}
}

func TestBuildReq(t *testing.T) {
	cl := NewClient("test-client-id", "test-client-secret")
	cl.Token = "test-token"

	req, _ := http.NewRequest("GET", cl.urlFor("/").String(), nil)
	req = cl.buildReq(req)

	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type should be application/json")
	}
	if req.Header.Get("Authorization") != "Bearer test-token" {
		t.Error("Authorization should be Bearer test-token but: ", req.Header.Get("Authorization"))
	}
}

func TestBuildReqAdditionalHeaders(t *testing.T) {
	cl := NewClient("test-client-id", "test-client-secret")
	cl.Token = "test-token"
	cl.AdditionalHeaders.Set("X-Custom", "value")

	req, _ := http.NewRequest("GET", cl.urlFor("/").String(), nil)
	req = cl.buildReq(req)

	if req.Header.Get("X-Custom") != "value" {
		t.Error("X-Custom header should be set")
	}
}

func TestEnsureTokenCaching(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth2/token" {
			callCount++
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("test-client-id", "test-client-secret", ts.URL, false)
	ctx := context.Background()

	client.ensureToken(ctx)
	client.ensureToken(ctx)
	client.ensureToken(ctx)

	if callCount != 1 {
		t.Errorf("expected 1 token request, got %d", callCount)
	}
}

func TestEnsureTokenRefreshOnExpiry(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/oauth2/token" {
			callCount++
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("test-client-id", "test-client-secret", ts.URL, false)
	ctx := context.Background()

	client.ensureToken(ctx)
	// Simulate expired token
	client.TokenExpiry = time.Now().Add(-1 * time.Minute)
	client.ensureToken(ctx)

	if callCount != 2 {
		t.Errorf("expected 2 token requests, got %d", callCount)
	}
}

func TestEnsureTokenServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("test-client-id", "test-client-secret", ts.URL, false)
	err := client.ensureToken(context.Background())
	if err == nil {
		t.Error("expected error for unauthorized token request")
	}
}

func TestPostJSON(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		body, _ := io.ReadAll(req.Body)
		var data map[string]any
		json.Unmarshal(body, &data)
		if data["name"] != "test" {
			t.Errorf("expected name=test, got %v", data["name"])
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(`{"id":1,"name":"test"}`))
	})
	defer ts.Close()

	resp, err := client.PostJSON(context.Background(), "/test", map[string]any{"name": "test"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestPutJSON(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte(`{"id":1,"name":"updated"}`))
	})
	defer ts.Close()

	resp, err := client.PutJSON(context.Background(), "/test/1", map[string]any{"name": "updated"})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSetPagination(t *testing.T) {
	cl := NewClient("id", "secret")
	u := cl.urlFor("/test")
	setPagination(u, 10, 20)
	if u.Query().Get("offset") != "10" {
		t.Errorf("expected offset=10, got %s", u.Query().Get("offset"))
	}
	if u.Query().Get("limit") != "20" {
		t.Errorf("expected limit=20, got %s", u.Query().Get("limit"))
	}
}

func TestSetPaginationZeroValues(t *testing.T) {
	cl := NewClient("id", "secret")
	u := cl.urlFor("/test")
	setPagination(u, 0, 0)
	if u.Query().Get("offset") != "" {
		t.Error("offset should not be set for zero value")
	}
	if u.Query().Get("limit") != "" {
		t.Error("limit should not be set for zero value")
	}
}
