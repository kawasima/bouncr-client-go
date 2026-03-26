package bouncr

import (
	"context"
	"net/http"
	"testing"
)

func TestCreatePasswordCredential(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/bouncr/api/password_credential" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"user":{"id":1,"account":"admin"},"created_at":"2026-01-01T00:00:00Z"}`))
	})
	defer ts.Close()

	cred, err := client.CreatePasswordCredential(context.Background(), &PasswordCredentialCreateRequest{
		Account:  "admin",
		Password: "secret",
		Initial:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if cred.User.Account != "admin" {
		t.Errorf("expected admin, got %s", cred.User.Account)
	}
}

func TestUpdatePasswordCredential(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Write([]byte(`{"user":{"id":1,"account":"admin"},"created_at":"2026-01-01T00:00:00Z"}`))
	})
	defer ts.Close()

	cred, err := client.UpdatePasswordCredential(context.Background(), &PasswordCredentialUpdateRequest{
		OldPassword: "old",
		NewPassword: "new",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cred.User.Account != "admin" {
		t.Errorf("expected admin, got %s", cred.User.Account)
	}
}

func TestDeletePasswordCredential(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeletePasswordCredential(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
