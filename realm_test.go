package bouncr

import (
	"context"
	"net/http"
	"testing"
)

func TestListRealms(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/application/myapp/realms" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`[{"id":1,"name":"default","url":""}]`))
	})
	defer ts.Close()

	realms, err := client.ListRealms(context.Background(), "myapp", &RealmSearchParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(realms) != 1 {
		t.Errorf("expected 1 realm, got %d", len(realms))
	}
}

func TestFindRealm(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/application/myapp/realm/default" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"name":"default","url":""}`))
	})
	defer ts.Close()

	realm, err := client.FindRealm(context.Background(), "myapp", "default")
	if err != nil {
		t.Fatal(err)
	}
	if realm.Name != "default" {
		t.Errorf("expected default, got %s", realm.Name)
	}
}

func TestCreateRealm(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		res.Write([]byte(`{"id":2,"name":"api","url":"api/"}`))
	})
	defer ts.Close()

	realm, err := client.CreateRealm(context.Background(), "myapp", &RealmCreateRequest{Name: "api", URL: "api/"})
	if err != nil {
		t.Fatal(err)
	}
	if realm.Name != "api" {
		t.Errorf("expected api, got %s", realm.Name)
	}
}

func TestUpdateRealm(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		res.Write([]byte(`{"id":1,"name":"default","description":"Updated"}`))
	})
	defer ts.Close()

	realm, err := client.UpdateRealm(context.Background(), "myapp", "default", &RealmUpdateRequest{Description: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if realm.Description != "Updated" {
		t.Errorf("expected Updated, got %s", realm.Description)
	}
}

func TestDeleteRealm(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeleteRealm(context.Background(), "myapp", "default")
	if err != nil {
		t.Fatal(err)
	}
}
