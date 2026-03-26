package bouncr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL  = "http://localhost:3000"
	defaultBasePath = "/bouncr/api"
	oauth2Path      = "/oauth2"
)

// Client is an API client for bouncr.
type Client struct {
	BaseURL           *url.URL
	TokenURL          *url.URL
	ClientID          string
	ClientSecret      string
	Token             string
	TokenExpiry       time.Time
	Verbose           bool
	AdditionalHeaders http.Header
	httpClient        *http.Client
}

// NewClient returns a new Client with default settings.
func NewClient(clientID, clientSecret string) *Client {
	u, _ := url.Parse(defaultBaseURL)
	return &Client{
		BaseURL:           u,
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		Verbose:           true,
		AdditionalHeaders: http.Header{},
		httpClient:        &http.Client{},
	}
}

// NewClientWithOptions returns a new Client with the given options.
func NewClientWithOptions(clientID, clientSecret, rawurl string, verbose bool) (*Client, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:           u,
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		Verbose:           verbose,
		AdditionalHeaders: http.Header{},
		httpClient:        &http.Client{},
	}, nil
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

func (c *Client) urlForOAuth2(path string) *url.URL {
	newURL, err := url.Parse(c.BaseURL.String())
	if err != nil {
		panic("invalid url passed")
	}

	newURL.Path = oauth2Path + path

	return newURL
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.Token != "" && time.Now().Before(c.TokenExpiry.Add(-10*time.Second)) {
		return nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	tokenURL := c.urlForOAuth2("/token").String()
	if c.TokenURL != nil {
		tokenURL = c.TokenURL.String()
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.ClientID+":"+c.ClientSecret)))

	if c.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			log.Printf("%s\n", dump)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			log.Printf("%s\n", dump)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed: %s", resp.Status)
	}

	var tokenResp tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return err
	}

	c.Token = tokenResp.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// Request sends a request to the bouncr API.
func (c *Client) Request(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}
	req = c.buildReq(req)

	if c.Verbose {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			log.Printf("%s\n", dump)
		}
	}
	resp, err = c.httpClient.Do(req)
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

// PostJSON is a shortcut method for posting json.
func (c *Client) PostJSON(ctx context.Context, path string, payload any) (*http.Response, error) {
	return c.requestJSON(ctx, "POST", path, payload)
}

// PutJSON is a shortcut method for putting json.
func (c *Client) PutJSON(ctx context.Context, path string, payload any) (*http.Response, error) {
	return c.requestJSON(ctx, "PUT", path, payload)
}

func (c *Client) requestJSON(ctx context.Context, method, path string, payload any) (*http.Response, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.urlFor(path).String(), &body)
	if err != nil {
		return nil, err
	}
	return c.Request(ctx, req)
}

func setPagination(u *url.URL, offset, limit int) {
	q := u.Query()
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	u.RawQuery = q.Encode()
}

func closeResponse(resp *http.Response) {
	if resp != nil {
		resp.Body.Close()
	}
}
