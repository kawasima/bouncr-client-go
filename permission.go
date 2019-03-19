package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Permission permission information
type Permission struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PermissionSearchParams parameters for search permissions
type PermissionSearchParams struct {
	Offset int
	Limit  int
}

// PermissionCreateRequest request for creating an permission
type PermissionCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// PermissionUpdateRequest request for creating an permission
type PermissionUpdateRequest PermissionCreateRequest

// FindPermission find a permission
func (c *Client) FindPermission(name string) (*Permission, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/permission/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *Permission
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err

}

// ListPermissions find the permissions
func (c *Client) ListPermissions(param *PermissionSearchParams) ([]*Permission, error) {
	req, err := http.NewRequest("GET", c.urlFor("/permissions").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(Permission)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreatePermission create an permission
func (c *Client) CreatePermission(createRequest *PermissionCreateRequest) (*Permission, error) {
	resp, err := c.PostJSON("/permissions", createRequest)
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

// UpdatePermission update an permission
func (c *Client) UpdatePermission(name string, updateRequest *PermissionUpdateRequest) (*Permission, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/permission/%s", name), updateRequest)
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

func (c *Client) DeletePermission(name string) error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/permission/%s", name)).String(),
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
