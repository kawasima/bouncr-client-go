package bouncr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestListRoles(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`[{"id":1,"name":"admin","description":"Admin role"}]`))
	})
	defer ts.Close()

	roles, err := client.ListRoles(context.Background(), &RoleSearchParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 || roles[0].Name != "admin" {
		t.Errorf("unexpected roles: %v", roles)
	}
}

func TestFindRole(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("embed") != "(permissions)" {
			t.Error("expected embed=(permissions)")
		}
		res.Write([]byte(`{"id":1,"name":"admin","permissions":[{"id":1,"name":"user:read"}]}`))
	})
	defer ts.Close()

	role, err := client.FindRole(context.Background(), "admin")
	if err != nil {
		t.Fatal(err)
	}
	if role.Name != "admin" {
		t.Errorf("expected admin, got %s", role.Name)
	}
}

func TestCreateRole(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		res.Write([]byte(`{"id":2,"name":"editor"}`))
	})
	defer ts.Close()

	role, err := client.CreateRole(context.Background(), &RoleCreateRequest{Name: "editor"})
	if err != nil {
		t.Fatal(err)
	}
	if role.Name != "editor" {
		t.Errorf("expected editor, got %s", role.Name)
	}
}

func TestUpdateRole(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Write([]byte(`{"id":1,"name":"admin","description":"Updated"}`))
	})
	defer ts.Close()

	role, err := client.UpdateRole(context.Background(), "admin", &RoleUpdateRequest{Description: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if role.Description != "Updated" {
		t.Errorf("expected Updated, got %s", role.Description)
	}
}

func TestDeleteRole(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeleteRole(context.Background(), "admin")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindPermissionsInRole(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/role/admin/permissions" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"name":"admin","permissions":[{"id":1,"name":"user:read"}]}`))
	})
	defer ts.Close()

	role, err := client.FindPermissionsInRole(context.Background(), "admin")
	if err != nil {
		t.Fatal(err)
	}
	if role.Permissions == nil || len(*role.Permissions) != 1 {
		t.Error("expected 1 permission")
	}
}

func TestAddPermissionsToRole(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		body, _ := io.ReadAll(req.Body)
		var perms []string
		json.Unmarshal(body, &perms)
		res.Write(body)
	})
	defer ts.Close()

	perms := []string{"user:read"}
	result, err := client.AddPermissionsToRole(context.Background(), "admin", &perms)
	if err != nil {
		t.Fatal(err)
	}
	if len(*result) != 1 {
		t.Errorf("expected 1 permission, got %d", len(*result))
	}
}

func TestRemovePermissionsFromRole(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	perms := []string{"user:read"}
	err := client.RemovePermissionsFromRole(context.Background(), "admin", &perms)
	if err != nil {
		t.Fatal(err)
	}
}
