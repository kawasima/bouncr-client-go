package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FindUsersInGroup finds users in a group.
func (c *Client) FindUsersInGroup(ctx context.Context, name string) (*Group, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(fmt.Sprintf("/group/%s/users", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Group](resp)
}

// AddUsersToGroup adds users to a group.
func (c *Client) AddUsersToGroup(ctx context.Context, group string, createRequest *[]string) (*[]string, error) {
	resp, err := c.PostJSON(ctx, fmt.Sprintf("/group/%s/users", group), createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return createRequest, nil
	}
	var data []string
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// RemoveUsersFromGroup removes users from a group.
func (c *Client) RemoveUsersFromGroup(ctx context.Context, name string, deleteRequest *[]string) error {
	resp, err := c.requestJSON(ctx, "DELETE", fmt.Sprintf("/group/%s/users", name), deleteRequest)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
