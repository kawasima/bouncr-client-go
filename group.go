package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Group group information
type Group struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Users       *[]User `json:"users"`
}

// GroupSearchParams parameters for search groups
type GroupSearchParams struct {
	Offset int
	Limit  int
}

// GroupCreateRequest request for creating an group
type GroupCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// GroupUpdateRequest request for creating an group
type GroupUpdateRequest GroupCreateRequest

// FindGroup find a group
func (c *Client) FindGroup(name string) (*Group, error) {
	url := c.urlFor(fmt.Sprintf("/group/%s", name))
	q := url.Query()
	q.Set("embed", "(users)")
	url.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", url.String(), nil)
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

// ListGroups find the groups
func (c *Client) ListGroups(param *GroupSearchParams) ([]*Group, error) {
	req, err := http.NewRequest("GET", c.urlFor("/groups").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(Group)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateGroup create an group
func (c *Client) CreateGroup(createRequest *GroupCreateRequest) (*Group, error) {
	resp, err := c.PostJSON("/groups", createRequest)
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

// UpdateGroup update an group
func (c *Client) UpdateGroup(name string, updateRequest *GroupUpdateRequest) (*Group, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/group/%s", name), updateRequest)
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

func (c *Client) DeleteGroup(name string) error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/group/%s", name)).String(),
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
