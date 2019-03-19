package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FindPermissionsInRole find permissions in a roles
func (c *Client) FindPermissionsInRole(name string) (*Role, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/role/%s/permissions", name)).String(), nil)
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

// AddPermissionsToRole add permissions to the role
func (c *Client) AddPermissionsToRole(role string, createRequest *[]string) (*[]string, error) {
	resp, err := c.PostJSON(fmt.Sprintf("/role/%s/permissions", role), createRequest)
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

func (c *Client) RemovePermissionsFromRole(name string, deleteRequest *[]string) error {
	resp, err := c.requestJSON("DELETE",
		c.urlFor(fmt.Sprintf("/role/%s/permissions", name)).String(),
		deleteRequest)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
