package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Permission represents permission information.
type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PermissionSearchParams contains parameters for searching permissions.
type PermissionSearchParams struct {
	Offset int
	Limit  int
}

// PermissionCreateRequest is a request for creating a permission.
type PermissionCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// PermissionUpdateRequest is a request for updating a permission.
type PermissionUpdateRequest PermissionCreateRequest

// FindPermission finds a permission by name.
func (c *Client) FindPermission(ctx context.Context, name string) (*Permission, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(fmt.Sprintf("/permission/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Permission
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// ListPermissions finds permissions.
func (c *Client) ListPermissions(ctx context.Context, param *PermissionSearchParams) ([]*Permission, error) {
	u := c.urlFor("/permissions")
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

	var data []*Permission
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreatePermission creates a permission.
func (c *Client) CreatePermission(ctx context.Context, createRequest *PermissionCreateRequest) (*Permission, error) {
	resp, err := c.PostJSON(ctx, "/permissions", createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Permission
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdatePermission updates a permission.
func (c *Client) UpdatePermission(ctx context.Context, name string, updateRequest *PermissionUpdateRequest) (*Permission, error) {
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/permission/%s", name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Permission
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// DeletePermission deletes a permission.
func (c *Client) DeletePermission(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/permission/%s", name)).String(), nil)
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
