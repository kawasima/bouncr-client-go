package bouncr

import (
	"context"
	"net/http"
	"testing"
)

func TestListUsers(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`[{"id":1,"account":"admin"},{"id":2,"account":"user1"}]`))
	})
	defer ts.Close()

	users, err := client.ListUsers(context.Background(), &UserSearchParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestListUsersWithParams(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		if q.Get("q") != "admin" {
			t.Errorf("expected q=admin, got %s", q.Get("q"))
		}
		if q.Get("offset") != "10" {
			t.Errorf("expected offset=10, got %s", q.Get("offset"))
		}
		if q.Get("limit") != "5" {
			t.Errorf("expected limit=5, got %s", q.Get("limit"))
		}
		if q.Get("group_id") != "3" {
			t.Errorf("expected group_id=3, got %s", q.Get("group_id"))
		}
		if q.Get("embed") != "(groups)" {
			t.Errorf("expected embed=(groups), got %s", q.Get("embed"))
		}
		res.Write([]byte(`[]`))
	})
	defer ts.Close()

	_, err := client.ListUsers(context.Background(), &UserSearchParams{
		Query:   "admin",
		Offset:  10,
		Limit:   5,
		GroupID: 3,
		Embed:   "(groups)",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListUsersNilParams(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.RawQuery != "" {
			t.Errorf("expected no query params, got %s", req.URL.RawQuery)
		}
		res.Write([]byte(`[]`))
	})
	defer ts.Close()

	_, err := client.ListUsers(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindUser(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/bouncr/api/user/admin" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"account":"admin","email":"admin@example.com"}`))
	})
	defer ts.Close()

	user, err := client.FindUser(context.Background(), "admin")
	if err != nil {
		t.Fatal(err)
	}
	if user.Account != "admin" {
		t.Errorf("expected admin, got %s", user.Account)
	}
	if user.ID != 1 {
		t.Errorf("expected ID=1, got %d", user.ID)
	}
	if user.UserProfiles["email"] != "admin@example.com" {
		t.Errorf("expected email in profiles, got %v", user.UserProfiles)
	}
}

func TestCreateUser(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		res.Write([]byte(`{"id":3,"account":"newuser"}`))
	})
	defer ts.Close()

	req := UserCreateRequest{"account": "newuser", "email": "new@example.com"}
	user, err := client.CreateUser(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	if user.Account != "newuser" {
		t.Errorf("expected newuser, got %s", user.Account)
	}
}

func TestUpdateUser(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			t.Errorf("expected PUT, got %s", req.Method)
		}
		if req.URL.Path != "/bouncr/api/user/admin" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.Write([]byte(`{"id":1,"account":"admin"}`))
	})
	defer ts.Close()

	req := UserUpdateRequest{"account": "admin", "email": "updated@example.com"}
	user, err := client.UpdateUser(context.Background(), "admin", &req)
	if err != nil {
		t.Fatal(err)
	}
	if user.Account != "admin" {
		t.Errorf("expected admin, got %s", user.Account)
	}
}

func TestDeleteUser(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	err := client.DeleteUser(context.Background(), "admin")
	if err != nil {
		t.Fatal(err)
	}
}
