package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OidcApplication oidcApplication information
type OidcApplication struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	HomeURL     string   `json:"home_url"`
	CallbackURL string   `json:"callback_url"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// OidcApplicationSearchParams parameters for search oidcApplications
type OidcApplicationSearchParams struct {
	Offset int
	Limit  int
}

// OidcApplicationCreateRequest request for creating an oidcApplication
type OidcApplicationCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	HomeURL     string   `json:"home_url"`
	CallbackURL string   `json:"callback_url"`
	Permissions []string `json:"permissions,omitempty"`
}

// OidcApplicationUpdateRequest request for creating an oidcApplication
type OidcApplicationUpdateRequest OidcApplicationCreateRequest

// FindOidcApplication find a oidcApplication
func (c *Client) FindOidcApplication(name string) (*OidcApplication, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/oidc_application/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *OidcApplication
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err

}

// ListOidcApplications find the oidcApplications
func (c *Client) ListOidcApplications(param *OidcApplicationSearchParams) ([]*OidcApplication, error) {
	req, err := http.NewRequest("GET", c.urlFor("/oidc_applications").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(OidcApplication)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateOidcApplication create an oidcApplication
func (c *Client) CreateOidcApplication(createRequest *OidcApplicationCreateRequest) (*OidcApplication, error) {
	resp, err := c.PostJSON("/oidc_applications", createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var data OidcApplication

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateOidcApplication update an oidcApplication
func (c *Client) UpdateOidcApplication(name string, updateRequest *OidcApplicationUpdateRequest) (*OidcApplication, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/oidc_application/%s", name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data OidcApplication

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) DeleteOidcApplication(name string) error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/oidc_application/%s", name)).String(),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
