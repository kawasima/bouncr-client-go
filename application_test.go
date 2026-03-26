package bouncr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestListApplications(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/applications" {
			t.Error("request URL should be /bouncr/api/applications but: ", req.URL.Path)
		}
		if req.Header.Get("Authorization") == "" {
			t.Error("request is NOT authenticated")
		}

		respJSON, _ := json.Marshal([]map[string]any{
			{"id": 1, "name": "demo1"},
		})
		res.Write(respJSON)
	})
	defer ts.Close()

	applications, err := client.ListApplications(context.Background(), &ApplicationSearchParams{})
	if err != nil {
		t.Fatal("err should be nil but: ", err)
	}
	if applications[0].Name != "demo1" {
		t.Error("expected name demo1 but: ", applications[0].Name)
	}
}

func TestListApplicationsWithPagination(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("offset") != "10" {
			t.Errorf("expected offset=10, got %s", req.URL.Query().Get("offset"))
		}
		if req.URL.Query().Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", req.URL.Query().Get("limit"))
		}
		res.Write([]byte(`[]`))
	})
	defer ts.Close()

	_, err := client.ListApplications(context.Background(), &ApplicationSearchParams{Offset: 10, Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/application/myapp" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		if req.URL.Query().Get("embed") != "(realms)" {
			t.Error("expected embed=(realms)")
		}
		res.Write([]byte(`{"id":1,"name":"myapp","description":"My App","virtual_path":"/myapp","pass_to":"http://backend:8080"}`))
	})
	defer ts.Close()

	app, err := client.FindApplication(context.Background(), "myapp")
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "myapp" {
		t.Errorf("expected myapp, got %s", app.Name)
	}
	if app.VirtualPath != "/myapp" {
		t.Errorf("expected /myapp, got %s", app.VirtualPath)
	}
}

func TestCreateApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		body, _ := io.ReadAll(req.Body)
		var data map[string]any
		json.Unmarshal(body, &data)
		if data["name"] != "newapp" {
			t.Errorf("expected name=newapp, got %v", data["name"])
		}
		res.Write([]byte(`{"id":2,"name":"newapp"}`))
	})
	defer ts.Close()

	app, err := client.CreateApplication(context.Background(), &ApplicationCreateRequest{Name: "newapp"})
	if err != nil {
		t.Fatal(err)
	}
	if app.ID != 2 {
		t.Errorf("expected ID=2, got %d", app.ID)
	}
}

func TestUpdateApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		if req.URL.Path != "/bouncr/api/application/myapp" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"name":"myapp","description":"Updated"}`))
	})
	defer ts.Close()

	app, err := client.UpdateApplication(context.Background(), "myapp", &ApplicationUpdateRequest{Description: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if app.Description != "Updated" {
		t.Errorf("expected Updated, got %s", app.Description)
	}
}

func TestDeleteApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		if req.URL.Path != "/bouncr/api/application/myapp" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.WriteHeader(http.StatusNoContent)
	})
	defer ts.Close()

	err := client.DeleteApplication(context.Background(), "myapp")
	if err != nil {
		t.Fatal(err)
	}
}
