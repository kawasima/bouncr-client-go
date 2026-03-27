package bouncr

import (
	"context"
	"fmt"
	"net/http"
)

// Role represents role information.
type Role struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Permissions *[]Permission `json:"permissions"`
}

// RoleSearchParams contains parameters for searching roles.
type RoleSearchParams struct {
	Offset int
	Limit  int
}

// RoleCreateRequest is a request for creating a role.
type RoleCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// RoleUpdateRequest is a request for updating a role.
type RoleUpdateRequest RoleCreateRequest

// FindRole finds a role by name.
func (c *Client) FindRole(ctx context.Context, name string) (*Role, error) {
	u := c.urlFor(fmt.Sprintf("/role/%s", name))
	q := u.Query()
	q.Set("embed", "(permissions)")
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Role](resp)
}

// ListRoles finds roles.
func (c *Client) ListRoles(ctx context.Context, param *RoleSearchParams) ([]*Role, error) {
	u := c.urlFor("/roles")
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
	return decodeJSONSlice[Role](resp)
}

// CreateRole creates a role.
func (c *Client) CreateRole(ctx context.Context, createRequest *RoleCreateRequest) (*Role, error) {
	resp, err := c.PostJSON(ctx, "/roles", createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Role](resp)
}

// UpdateRole updates a role.
func (c *Client) UpdateRole(ctx context.Context, name string, updateRequest *RoleUpdateRequest) (*Role, error) {
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/role/%s", name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Role](resp)
}

// DeleteRole deletes a role.
func (c *Client) DeleteRole(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/role/%s", name)).String(), nil)
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
