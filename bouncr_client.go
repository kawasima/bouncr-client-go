package bouncr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	defaultBaseURL  = "http://localhost:3000"
	defaultBasePath = "/bouncr/api"
)

// Client api client for bouncr
type Client struct {
	BaseURL           *url.URL
	Account           string
	Password          string
	Token             string
	Verbose           bool
	AdditionalHeaders http.Header
}

// NewClient returns new mackerel.Client
func NewClient(account string, password string) *Client {
	u, _ := url.Parse(defaultBaseURL)
	return &Client{u, account, password, "", true, http.Header{}}
}

func NewClientWithOptions(account string, password string, rawurl string, verbose bool) (*Client, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Client{u, account, password, "", verbose, http.Header{}}, nil
}

func (c *Client) buildReq(req *http.Request) *http.Request {
	for header, values := range c.AdditionalHeaders {
		for _, v := range values {
			req.Header.Add(header, v)
		}
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req
}

func (c *Client) urlFor(path string) *url.URL {
	newURL, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("invalid url passed")
	}

	newURL.Path = defaultBasePath + path

	return newURL
}

// Request request to mackerel and receive response
func (c *Client) Request(req *http.Request) (resp *http.Response, err error) {
	if !strings.HasSuffix(req.URL.Path, "/sign_in") && c.Token == "" {
		signInResponse, err := c.SignIn(&SignInRequest{
			Account:  c.Account,
			Password: c.Password,
		})
		if err != nil {
			panic(err)
		}
		c.Token = signInResponse.Token
	}
	req = c.buildReq(req)

	if c.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			fmt.Printf("%s\n", dump)
		}
	}
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			log.Printf("%s\n", dump)
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("API result failed: %s", resp.Status)
	}
	return resp, nil
}

// PostJSON shortcut method for posting json
func (c *Client) PostJSON(path string, payload interface{}) (*http.Response, error) {
	return c.requestJSON("POST", path, payload)
}

// PutJSON shortcut method for putting json
func (c *Client) PutJSON(path string, payload interface{}) (*http.Response, error) {
	return c.requestJSON("PUT", path, payload)
}

func (c *Client) requestJSON(method string, path string, payload interface{}) (*http.Response, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.urlFor(path).String(), &body)
	if err != nil {
		return nil, err
	}
	return c.Request(req)
}
func closeResponse(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
