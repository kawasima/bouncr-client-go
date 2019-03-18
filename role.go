package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Role role information
type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RoleSearchParams parameters for search roles
type RoleSearchParams struct {
	Offset int
	Limit  int
}

// RoleCreateRequest request for creating an role
type RoleCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// RoleUpdateRequest request for creating an role
type RoleUpdateRequest RoleCreateRequest

// FindRole find a role
func (c *Client) FindRole(name string) (*Role, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/role/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *Role
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err

}

// ListRoles find the roles
func (c *Client) ListRoles(param *RoleSearchParams) ([]*Role, error) {
	req, err := http.NewRequest("GET", c.urlFor("/roles").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(Role)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateRole create an role
func (c *Client) CreateRole(createRequest *RoleCreateRequest) (*Role, error) {
	resp, err := c.PostJSON("/roles", createRequest)
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

// UpdateRole update an role
func (c *Client) UpdateRole(name string, updateRequest *RoleUpdateRequest) (*Role, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/role/%s", name), updateRequest)
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
