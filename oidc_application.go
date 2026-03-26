package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// OidcApplication represents OIDC application information.
type OidcApplication struct {
	ID                    int          `json:"id"`
	Name                  string       `json:"name"`
	Description           string       `json:"description"`
	ClientID              string       `json:"client_id,omitempty"`
	HomeURI               string       `json:"home_uri,omitempty"`
	CallbackURI           string       `json:"callback_uri,omitempty"`
	BackchannelLogoutURI  string       `json:"backchannel_logout_uri,omitempty"`
	FrontchannelLogoutURI string       `json:"frontchannel_logout_uri,omitempty"`
	GrantTypes            []string     `json:"grant_types,omitempty"`
	Permissions           []Permission `json:"permissions,omitempty"`
}

// OidcApplicationSearchParams contains parameters for searching OIDC applications.
type OidcApplicationSearchParams struct {
	Offset int
	Limit  int
}

// OidcApplicationCreateRequest is a request for creating an OIDC application.
type OidcApplicationCreateRequest struct {
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	HomeURI               string   `json:"home_uri,omitempty"`
	CallbackURI           string   `json:"callback_uri,omitempty"`
	BackchannelLogoutURI  string   `json:"backchannel_logout_uri,omitempty"`
	FrontchannelLogoutURI string   `json:"frontchannel_logout_uri,omitempty"`
	GrantTypes            []string `json:"grant_types,omitempty"`
	Permissions           []string `json:"permissions,omitempty"`
}

// OidcApplicationUpdateRequest is a request for updating an OIDC application.
type OidcApplicationUpdateRequest OidcApplicationCreateRequest

// OidcApplicationSecret represents a client_id and client_secret pair.
type OidcApplicationSecret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// FindOidcApplication finds an OIDC application by name.
func (c *Client) FindOidcApplication(ctx context.Context, name string) (*OidcApplication, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(fmt.Sprintf("/oidc_application/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
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

// ListOidcApplications finds OIDC applications.
func (c *Client) ListOidcApplications(ctx context.Context, param *OidcApplicationSearchParams) ([]*OidcApplication, error) {
	u := c.urlFor("/oidc_applications")
	if param != nil {
		setPagination(u, param.Offset, param.Limit)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*OidcApplication
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateOidcApplication creates an OIDC application.
func (c *Client) CreateOidcApplication(ctx context.Context, createRequest *OidcApplicationCreateRequest) (*OidcApplication, error) {
	resp, err := c.PostJSON(ctx, "/oidc_applications", createRequest)
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

// UpdateOidcApplication updates an OIDC application.
func (c *Client) UpdateOidcApplication(ctx context.Context, name string, updateRequest *OidcApplicationUpdateRequest) (*OidcApplication, error) {
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/oidc_application/%s", name), updateRequest)
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

// DeleteOidcApplication deletes an OIDC application.
func (c *Client) DeleteOidcApplication(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/oidc_application/%s", name)).String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}

// RegenerateOidcApplicationSecret regenerates the client secret for an OIDC application.
func (c *Client) RegenerateOidcApplicationSecret(ctx context.Context, name string) (*OidcApplicationSecret, error) {
	resp, err := c.PostJSON(ctx, fmt.Sprintf("/oidc_application/%s/secret", name), nil)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data OidcApplicationSecret
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
