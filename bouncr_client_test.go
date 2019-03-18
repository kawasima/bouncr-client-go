package bouncr

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

	}))
	defer ts.Close()
	client := NewClient("admin", "password")

	req, _ := http.NewRequest("GET", client.urlFor("./").String(), nil)
	client.Request(req)
}

func TestBuildReq(t *testing.T) {
	cl := NewClient("admin", "password")

	req, _ := http.NewRequest("GET", cl.urlFor("/").String(), nil)
	req = cl.buildReq(req)

	log.Println(req)
}
