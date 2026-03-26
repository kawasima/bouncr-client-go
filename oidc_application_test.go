package bouncr

import (
	"context"
	"net/http"
	"testing"
)

func TestListOidcApplications(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`[{"id":1,"name":"myoidc","client_id":"abc"}]`))
	})
	defer ts.Close()

	apps, err := client.ListOidcApplications(context.Background(), &OidcApplicationSearchParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 1 || apps[0].ClientID != "abc" {
		t.Errorf("unexpected apps: %v", apps)
	}
}

func TestFindOidcApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/oidc_application/myoidc" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"name":"myoidc","client_id":"abc","grant_types":["client_credentials"]}`))
	})
	defer ts.Close()

	app, err := client.FindOidcApplication(context.Background(), "myoidc")
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "myoidc" {
		t.Errorf("expected myoidc, got %s", app.Name)
	}
	if len(app.GrantTypes) != 1 {
		t.Errorf("expected 1 grant type, got %d", len(app.GrantTypes))
	}
}

func TestCreateOidcApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		res.Write([]byte(`{"id":2,"name":"newoidc","client_id":"xyz"}`))
	})
	defer ts.Close()

	app, err := client.CreateOidcApplication(context.Background(), &OidcApplicationCreateRequest{Name: "newoidc"})
	if err != nil {
		t.Fatal(err)
	}
	if app.ClientID != "xyz" {
		t.Errorf("expected xyz, got %s", app.ClientID)
	}
}

func TestUpdateOidcApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Write([]byte(`{"id":1,"name":"myoidc","description":"Updated"}`))
	})
	defer ts.Close()

	app, err := client.UpdateOidcApplication(context.Background(), "myoidc", &OidcApplicationUpdateRequest{Description: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if app.Description != "Updated" {
		t.Errorf("expected Updated, got %s", app.Description)
	}
}

func TestDeleteOidcApplication(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeleteOidcApplication(context.Background(), "myoidc")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegenerateOidcApplicationSecret(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/bouncr/api/oidc_application/myoidc/secret" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"client_id":"abc","client_secret":"newsecret"}`))
	})
	defer ts.Close()

	secret, err := client.RegenerateOidcApplicationSecret(context.Background(), "myoidc")
	if err != nil {
		t.Fatal(err)
	}
	if secret.ClientSecret != "newsecret" {
		t.Errorf("expected newsecret, got %s", secret.ClientSecret)
	}
}
