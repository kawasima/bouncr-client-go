package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Group represents group information.
type Group struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Users       *[]User `json:"users"`
}

// GroupSearchParams contains parameters for searching groups.
type GroupSearchParams struct {
	Offset int
	Limit  int
}

// GroupCreateRequest is a request for creating a group.
type GroupCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// GroupUpdateRequest is a request for updating a group.
type GroupUpdateRequest GroupCreateRequest

// FindGroup finds a group by name.
func (c *Client) FindGroup(ctx context.Context, name string) (*Group, error) {
	u := c.urlFor(fmt.Sprintf("/group/%s", name))
	q := u.Query()
	q.Set("embed", "(users)")
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

	var data Group
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// ListGroups finds groups.
func (c *Client) ListGroups(ctx context.Context, param *GroupSearchParams) ([]*Group, error) {
	u := c.urlFor("/groups")
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

	var data []*Group
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateGroup creates a group.
func (c *Client) CreateGroup(ctx context.Context, createRequest *GroupCreateRequest) (*Group, error) {
	resp, err := c.PostJSON(ctx, "/groups", createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Group
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateGroup updates a group.
func (c *Client) UpdateGroup(ctx context.Context, name string, updateRequest *GroupUpdateRequest) (*Group, error) {
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/group/%s", name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Group
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// DeleteGroup deletes a group.
func (c *Client) DeleteGroup(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/group/%s", name)).String(), nil)
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
