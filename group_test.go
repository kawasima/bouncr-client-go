package bouncr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestListGroups(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/groups" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`[{"id":1,"name":"admins","description":"Admins"}]`))
	})
	defer ts.Close()

	groups, err := client.ListGroups(context.Background(), &GroupSearchParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 || groups[0].Name != "admins" {
		t.Errorf("unexpected groups: %v", groups)
	}
}

func TestListGroupsWithPagination(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("offset") != "5" {
			t.Errorf("expected offset=5, got %s", req.URL.Query().Get("offset"))
		}
		res.Write([]byte(`[]`))
	})
	defer ts.Close()

	_, err := client.ListGroups(context.Background(), &GroupSearchParams{Offset: 5})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindGroup(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/group/admins" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		if req.URL.Query().Get("embed") != "(users)" {
			t.Error("expected embed=(users)")
		}
		res.Write([]byte(`{"id":1,"name":"admins","description":"Admins"}`))
	})
	defer ts.Close()

	group, err := client.FindGroup(context.Background(), "admins")
	if err != nil {
		t.Fatal(err)
	}
	if group.Name != "admins" {
		t.Errorf("expected admins, got %s", group.Name)
	}
}

func TestCreateGroup(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		res.Write([]byte(`{"id":2,"name":"newgroup"}`))
	})
	defer ts.Close()

	group, err := client.CreateGroup(context.Background(), &GroupCreateRequest{Name: "newgroup"})
	if err != nil {
		t.Fatal(err)
	}
	if group.Name != "newgroup" {
		t.Errorf("expected newgroup, got %s", group.Name)
	}
}

func TestUpdateGroup(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Write([]byte(`{"id":1,"name":"admins","description":"Updated"}`))
	})
	defer ts.Close()

	group, err := client.UpdateGroup(context.Background(), "admins", &GroupUpdateRequest{Description: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if group.Description != "Updated" {
		t.Errorf("expected Updated, got %s", group.Description)
	}
}

func TestDeleteGroup(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeleteGroup(context.Background(), "admins")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindUsersInGroup(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/group/admins/users" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"name":"admins","users":[{"id":1,"account":"admin"}]}`))
	})
	defer ts.Close()

	group, err := client.FindUsersInGroup(context.Background(), "admins")
	if err != nil {
		t.Fatal(err)
	}
	if group.Users == nil || len(*group.Users) != 1 {
		t.Error("expected 1 user")
	}
}

func TestAddUsersToGroup(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		body, _ := io.ReadAll(req.Body)
		var accounts []string
		json.Unmarshal(body, &accounts)
		if len(accounts) != 1 || accounts[0] != "user1" {
			t.Errorf("unexpected body: %s", body)
		}
		res.Write([]byte(`["user1"]`))
	})
	defer ts.Close()

	accounts := []string{"user1"}
	result, err := client.AddUsersToGroup(context.Background(), "admins", &accounts)
	if err != nil {
		t.Fatal(err)
	}
	if len(*result) != 1 {
		t.Errorf("expected 1 account, got %d", len(*result))
	}
}

func TestRemoveUsersFromGroup(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	accounts := []string{"user1"}
	err := client.RemoveUsersFromGroup(context.Background(), "admins", &accounts)
	if err != nil {
		t.Fatal(err)
	}
}
