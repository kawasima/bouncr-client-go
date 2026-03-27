package bouncr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// User represents user information.
type User struct {
	ID           int64          `json:"id"`
	Account      string         `json:"account"`
	UserProfiles map[string]any `json:"user_profiles"`
}

// UserSearchParams contains parameters for searching users.
type UserSearchParams struct {
	Query   string
	Offset  int
	Limit   int
	GroupID int
	Embed   string
}

// UserCreateRequest is a request for creating a user.
type UserCreateRequest map[string]any

// UserUpdateRequest is a request for updating a user.
type UserUpdateRequest map[string]any

// FindUser finds a user by account.
func (c *Client) FindUser(ctx context.Context, account string) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(fmt.Sprintf("/user/%s", account)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(ctx, req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return &User{}, nil
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	user := &User{
		ID:      int64(data["id"].(float64)),
		Account: data["account"].(string),
	}
	delete(data, "id")
	delete(data, "account")
	user.UserProfiles = data
	return user, nil
}

// ListUsers finds users.
func (c *Client) ListUsers(ctx context.Context, param *UserSearchParams) ([]*User, error) {
	u := c.urlFor("/users")
	if param != nil {
		setPagination(u, param.Offset, param.Limit)
		q := u.Query()
		if param.Query != "" {
			q.Set("q", param.Query)
		}
		if param.GroupID > 0 {
			q.Set("group_id", strconv.Itoa(param.GroupID))
		}
		if param.Embed != "" {
			q.Set("embed", param.Embed)
		}
		u.RawQuery = q.Encode()
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
	return decodeJSONSlice[User](resp)
}

// CreateUser creates a user.
func (c *Client) CreateUser(ctx context.Context, createRequest *UserCreateRequest) (*User, error) {
	resp, err := c.PostJSON(ctx, "/users", createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[User](resp)
}

// UpdateUser updates a user.
func (c *Client) UpdateUser(ctx context.Context, account string, updateRequest *UserUpdateRequest) (*User, error) {
	delete(*updateRequest, "account")
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/user/%s", account), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[User](resp)
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/user/%s", name)).String(), nil)
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
