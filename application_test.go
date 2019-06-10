package bouncr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListApplications(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/bouncr/api/sign_in" {
			respJSON, _ := json.Marshal(map[string]interface{}{
				"token": "550e8400-e29b-41d4-a716-446655440000",
			})

			fmt.Fprint(res, string(respJSON))
			return
		}

		if req.URL.Path != "/bouncr/api/applications" {
			t.Error("request URL should be /bouncr/api/applications but: ", req.URL.Path)
		}

		if req.Header.Get("Authorization") != "" {
			t.Error("request is NOT authenticated ")
		}

		respJSON, _ := json.Marshal([]map[string]interface{}{
			{
				"id":   1,
				"name": "demo1",
			},
		})

		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	log.Println(ts.URL)

	client, _ := NewClientWithOptions("admin", "password", ts.URL, false)
	applications, err := client.ListApplications(&ApplicationSearchParams{})

	log.Println(applications)
	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if applications[0].Name != "demo1" {
		t.Error("request sends json including name but: ", applications[0])
	}
}
