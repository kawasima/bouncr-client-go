package bouncr

import (
	"context"
	"net/http"
	"testing"
)

func TestListPermissions(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`[{"id":1,"name":"user:read"},{"id":2,"name":"user:create"}]`))
	})
	defer ts.Close()

	perms, err := client.ListPermissions(context.Background(), &PermissionSearchParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(perms) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(perms))
	}
}

func TestFindPermission(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/permission/user:read" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"name":"user:read","description":"Read users"}`))
	})
	defer ts.Close()

	perm, err := client.FindPermission(context.Background(), "user:read")
	if err != nil {
		t.Fatal(err)
	}
	if perm.Name != "user:read" {
		t.Errorf("expected user:read, got %s", perm.Name)
	}
}

func TestCreatePermission(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		res.Write([]byte(`{"id":3,"name":"user:delete"}`))
	})
	defer ts.Close()

	perm, err := client.CreatePermission(context.Background(), &PermissionCreateRequest{Name: "user:delete"})
	if err != nil {
		t.Fatal(err)
	}
	if perm.Name != "user:delete" {
		t.Errorf("expected user:delete, got %s", perm.Name)
	}
}

func TestUpdatePermission(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Write([]byte(`{"id":1,"name":"user:read","description":"Updated"}`))
	})
	defer ts.Close()

	perm, err := client.UpdatePermission(context.Background(), "user:read", &PermissionUpdateRequest{Description: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if perm.Description != "Updated" {
		t.Errorf("expected Updated, got %s", perm.Description)
	}
}

func TestDeletePermission(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeletePermission(context.Background(), "user:read")
	if err != nil {
		t.Fatal(err)
	}
}
