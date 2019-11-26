package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OidcProvider oidcProvider information
type OidcProvider struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	ClientId                string `json:"client_id"`
	Scope                   string `json:"scope"`
	ResponseType            string `json:"response_type"`
	AuthorizationEndpoint   string `json:"authorization_endpoint"`
	TokenEndpoint           string `json:"token_endpoint"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
	RedirectUri             string `json:"redirect_uri"`
}

// OidcProviderSearchParams parameters for search oidcProviders
type OidcProviderSearchParams struct {
	Offset int
	Limit  int
}

// OidcProviderCreateRequest request for creating an oidcProvider
type OidcProviderCreateRequest struct {
	Name                    string `json:"name"`
	ClientId                string `json:"client_id"`
	ClientSecret            string `json:"client_secret"`
	Scope                   string `json:"scope"`
	ResponseType            string `json:"response_type"`
	AuthorizationEndpoint   string `json:"authorization_endpoint"`
	TokenEndpoint           string `json:"token_endpoint"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
	RedirectUri             string `json:"redirect_uri"`
}

// OidcProviderUpdateRequest request for creating an oidcProvider
type OidcProviderUpdateRequest OidcProviderCreateRequest

// FindOidcProvider find a oidcProvider
func (c *Client) FindOidcProvider(name string) (*OidcProvider, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/oidc_provider/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *OidcProvider
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err

}

// ListOidcProviders find the oidcProviders
func (c *Client) ListOidcProviders(param *OidcProviderSearchParams) ([]*OidcProvider, error) {
	req, err := http.NewRequest("GET", c.urlFor("/oidc_providers").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(OidcProvider)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateOidcProvider create an oidcProvider
func (c *Client) CreateOidcProvider(createRequest *OidcProviderCreateRequest) (*OidcProvider, error) {
	resp, err := c.PostJSON("/oidc_providers", createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var data OidcProvider

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateOidcProvider update an oidcProvider
func (c *Client) UpdateOidcProvider(name string, updateRequest *OidcProviderUpdateRequest) (*OidcProvider, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/oidc_provider/%s", name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data OidcProvider

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) DeleteOidcProvider(name string) error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/oidc_provider/%s", name)).String(),
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
