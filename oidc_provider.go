package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// OidcProvider represents OIDC provider information.
type OidcProvider struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	ClientID                string `json:"client_id"`
	ClientSecret            string `json:"client_secret"`
	Scope                   string `json:"scope"`
	ResponseType            string `json:"response_type"`
	AuthorizationEndpoint   string `json:"authorization_endpoint"`
	TokenEndpoint           string `json:"token_endpoint"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
	RedirectURI             string `json:"redirect_uri"`
	PkceEnabled             bool   `json:"pkce_enabled"`
	JwksURI                 string `json:"jwks_uri,omitempty"`
	Issuer                  string `json:"issuer,omitempty"`
}

// OidcProviderSearchParams contains parameters for searching OIDC providers.
type OidcProviderSearchParams struct {
	Offset int
	Limit  int
}

// OidcProviderCreateRequest is a request for creating an OIDC provider.
type OidcProviderCreateRequest struct {
	Name                    string `json:"name"`
	ClientID                string `json:"client_id"`
	ClientSecret            string `json:"client_secret"`
	Scope                   string `json:"scope"`
	ResponseType            string `json:"response_type"`
	AuthorizationEndpoint   string `json:"authorization_endpoint"`
	TokenEndpoint           string `json:"token_endpoint"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
	RedirectURI             string `json:"redirect_uri"`
	PkceEnabled             bool   `json:"pkce_enabled"`
	JwksURI                 string `json:"jwks_uri,omitempty"`
	Issuer                  string `json:"issuer,omitempty"`
}

// OidcProviderUpdateRequest is a request for updating an OIDC provider.
type OidcProviderUpdateRequest OidcProviderCreateRequest

// FindOidcProvider finds an OIDC provider by name.
func (c *Client) FindOidcProvider(ctx context.Context, name string) (*OidcProvider, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(fmt.Sprintf("/oidc_provider/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
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

// ListOidcProviders finds OIDC providers.
func (c *Client) ListOidcProviders(ctx context.Context, param *OidcProviderSearchParams) ([]*OidcProvider, error) {
	u := c.urlFor("/oidc_providers")
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

	var data []*OidcProvider
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateOidcProvider creates an OIDC provider.
func (c *Client) CreateOidcProvider(ctx context.Context, createRequest *OidcProviderCreateRequest) (*OidcProvider, error) {
	resp, err := c.PostJSON(ctx, "/oidc_providers", createRequest)
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

// UpdateOidcProvider updates an OIDC provider.
func (c *Client) UpdateOidcProvider(ctx context.Context, name string, updateRequest *OidcProviderUpdateRequest) (*OidcProvider, error) {
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/oidc_provider/%s", name), updateRequest)
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

// DeleteOidcProvider deletes an OIDC provider.
func (c *Client) DeleteOidcProvider(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/oidc_provider/%s", name)).String(), nil)
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
