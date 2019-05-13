package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User user information
type User struct {
	ID           int64                    `json:"id"`
	Account      string                 `json:"account"`
	UserProfiles map[string]interface{} `json:"user_profiles"`
}

// UserSearchParams parameters for search users
type UserSearchParams struct {
	Offset int
	Limit  int
}

// UserCreateRequest request for creating an user
type UserCreateRequest map[string]interface{}
// UserUpdateRequest request for creating an user
type UserUpdateRequest map[string]interface{}

// FindUser find a user
func (c *Client) FindUser(account string) (*User, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/user/%s", account)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:           int64((*data)["id"].(float64)),
		Account:      (*data)["account"].(string),
	}
	delete(*data, "id")
	delete(*data, "account")
	user.UserProfiles = *data
	return user, err

}

// ListUsers find the users
func (c *Client) ListUsers(param *UserSearchParams) ([]*User, error) {
	req, err := http.NewRequest("GET", c.urlFor("/users").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(User)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateUser create an user
func (c *Client) CreateUser(createRequest *UserCreateRequest) (*User, error) {
	resp, err := c.PostJSON("/users", createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var data User

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateUser update an user
func (c *Client) UpdateUser(account string, updateRequest *UserUpdateRequest) (*User, error) {
	delete(*updateRequest, "account");
	resp, err := c.PutJSON(fmt.Sprintf("/user/%s", account), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data User
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) DeleteUser(name string) error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor(fmt.Sprintf("/user/%s", name)).String(),
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
