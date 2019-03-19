package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FindUsersInGroup find users in a groups
func (c *Client) FindUsersInGroup(name string) (*Group, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/group/%s/users", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *Group
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err

}

// AddUsersToGroup add users to the group
func (c *Client) AddUsersToGroup(group string, createRequest *[]string) (*[]string, error) {
	resp, err := c.PostJSON(fmt.Sprintf("/group/%/users", group), createRequest)
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

func (c *Client) RemoveUsersFromGroup(name string, deleteRequest *[]string) error {
	resp, err := c.requestJSON("DELETE",
		c.urlFor(fmt.Sprintf("/group/%s/users", name)).String(),
		deleteRequest)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
