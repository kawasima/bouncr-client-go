package bouncr

import (
	"context"
	"net/http"
	"testing"
)

func TestListOidcProviders(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`[{"id":1,"name":"google","client_id":"gid"}]`))
	})
	defer ts.Close()

	providers, err := client.ListOidcProviders(context.Background(), &OidcProviderSearchParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(providers) != 1 || providers[0].Name != "google" {
		t.Errorf("unexpected providers: %v", providers)
	}
}

func TestFindOidcProvider(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/oidc_provider/google" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"name":"google","client_id":"gid","authorization_endpoint":"https://accounts.google.com/o/oauth2/auth"}`))
	})
	defer ts.Close()

	provider, err := client.FindOidcProvider(context.Background(), "google")
	if err != nil {
		t.Fatal(err)
	}
	if provider.Name != "google" {
		t.Errorf("expected google, got %s", provider.Name)
	}
}

func TestCreateOidcProvider(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		res.Write([]byte(`{"id":2,"name":"github","client_id":"ghid"}`))
	})
	defer ts.Close()

	provider, err := client.CreateOidcProvider(context.Background(), &OidcProviderCreateRequest{Name: "github", ClientID: "ghid"})
	if err != nil {
		t.Fatal(err)
	}
	if provider.Name != "github" {
		t.Errorf("expected github, got %s", provider.Name)
	}
}

func TestUpdateOidcProvider(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Write([]byte(`{"id":1,"name":"google","scope":"openid email"}`))
	})
	defer ts.Close()

	provider, err := client.UpdateOidcProvider(context.Background(), "google", &OidcProviderUpdateRequest{Scope: "openid email"})
	if err != nil {
		t.Fatal(err)
	}
	if provider.Scope != "openid email" {
		t.Errorf("expected openid email, got %s", provider.Scope)
	}
}

func TestDeleteOidcProvider(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeleteOidcProvider(context.Background(), "google")
	if err != nil {
		t.Fatal(err)
	}
}
