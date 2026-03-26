package bouncr

import (
	"context"
	"net/http"
	"testing"
)

func TestFindAssignment(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		if q.Get("group") != "admins" {
			t.Errorf("expected group=admins, got %s", q.Get("group"))
		}
		if q.Get("role") != "admin" {
			t.Errorf("expected role=admin, got %s", q.Get("role"))
		}
		if q.Get("realm") != "default" {
			t.Errorf("expected realm=default, got %s", q.Get("realm"))
		}
		res.Write([]byte(`{"group":{"id":1,"name":"admins"},"role":{"id":1,"name":"admin"},"realm":{"id":1,"name":"default"}}`))
	})
	defer ts.Close()

	assignment, err := client.FindAssignment(context.Background(), &AssignmentRequest{
		Group: IDObject{Name: "admins"},
		Role:  IDObject{Name: "admin"},
		Realm: IDObject{Name: "default"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if assignment.Group.Name != "admins" {
		t.Errorf("expected admins, got %s", assignment.Group.Name)
	}
}

func TestCreateAssignments(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/bouncr/api/assignments" {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		res.WriteHeader(200)
	})
	defer ts.Close()

	reqs := []AssignmentRequest{
		{
			Group: IDObject{Name: "admins"},
			Role:  IDObject{Name: "admin"},
			Realm: IDObject{Name: "default"},
		},
	}
	_, err := client.CreateAssignments(context.Background(), &reqs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAssignments(t *testing.T) {
	ts, client := newTestServer(t, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", req.Method)
		}
		res.WriteHeader(204)
	})
	defer ts.Close()

	reqs := []AssignmentRequest{
		{
			Group: IDObject{Name: "admins"},
			Role:  IDObject{Name: "admin"},
			Realm: IDObject{Name: "default"},
		},
	}
	err := client.DeleteAssignments(context.Background(), &reqs)
	if err != nil {
		t.Fatal(err)
	}
}
