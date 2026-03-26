package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// FindPermissionsInRole finds permissions in a role.
func (c *Client) FindPermissionsInRole(ctx context.Context, name string) (*Role, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(fmt.Sprintf("/role/%s/permissions", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Role
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// AddPermissionsToRole adds permissions to a role.
func (c *Client) AddPermissionsToRole(ctx context.Context, role string, createRequest *[]string) (*[]string, error) {
	resp, err := c.PostJSON(ctx, fmt.Sprintf("/role/%s/permissions", role), createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// RemovePermissionsFromRole removes permissions from a role.
func (c *Client) RemovePermissionsFromRole(ctx context.Context, name string, deleteRequest *[]string) error {
	resp, err := c.requestJSON(ctx, "DELETE", fmt.Sprintf("/role/%s/permissions", name), deleteRequest)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
