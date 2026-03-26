package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Realm represents realm information.
type Realm struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// RealmSearchParams contains parameters for searching realms.
type RealmSearchParams struct {
	Offset int
	Limit  int
}

// RealmCreateRequest is a request for creating a realm.
type RealmCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// RealmUpdateRequest is a request for updating a realm.
type RealmUpdateRequest RealmCreateRequest

// FindRealm finds a realm by name within an application.
func (c *Client) FindRealm(ctx context.Context, applicationName, name string) (*Realm, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(fmt.Sprintf("/application/%s/realm/%s", applicationName, name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Realm
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// ListRealms finds realms within an application.
func (c *Client) ListRealms(ctx context.Context, applicationName string, param *RealmSearchParams) ([]*Realm, error) {
	u := c.urlFor(fmt.Sprintf("/application/%s/realms", applicationName))
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

	var data []*Realm
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateRealm creates a realm within an application.
func (c *Client) CreateRealm(ctx context.Context, applicationName string, createRequest *RealmCreateRequest) (*Realm, error) {
	resp, err := c.PostJSON(ctx, fmt.Sprintf("/application/%s/realms", applicationName), createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Realm
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateRealm updates a realm within an application.
func (c *Client) UpdateRealm(ctx context.Context, applicationName, name string, updateRequest *RealmUpdateRequest) (*Realm, error) {
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/application/%s/realm/%s", applicationName, name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Realm
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// DeleteRealm deletes a realm within an application.
func (c *Client) DeleteRealm(ctx context.Context, applicationName, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/application/%s/realm/%s", applicationName, name)).String(), nil)
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
