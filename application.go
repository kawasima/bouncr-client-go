package bouncr

import (
	"context"
	"fmt"
	"net/http"
)

// Application represents application information.
type Application struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PassTo      string `json:"pass_to"`
	VirtualPath string `json:"virtual_path"`
	TopPage     string `json:"top_page"`
}

// ApplicationSearchParams contains parameters for searching applications.
type ApplicationSearchParams struct {
	Offset int
	Limit  int
}

// ApplicationCreateRequest is a request for creating an application.
type ApplicationCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PassTo      string `json:"pass_to,omitempty"`
	VirtualPath string `json:"virtual_path,omitempty"`
	TopPage     string `json:"top_page,omitempty"`
}

// ApplicationUpdateRequest is a request for updating an application.
type ApplicationUpdateRequest ApplicationCreateRequest

// FindApplication finds an application by name.
func (c *Client) FindApplication(ctx context.Context, name string) (*Application, error) {
	u := c.urlFor(fmt.Sprintf("/application/%s", name))
	q := u.Query()
	q.Set("embed", "(realms)")
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
	return decodeJSON[Application](resp)
}

// ListApplications finds applications.
func (c *Client) ListApplications(ctx context.Context, param *ApplicationSearchParams) ([]*Application, error) {
	u := c.urlFor("/applications")
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
	return decodeJSONSlice[Application](resp)
}

// CreateApplication creates an application.
func (c *Client) CreateApplication(ctx context.Context, createRequest *ApplicationCreateRequest) (*Application, error) {
	resp, err := c.PostJSON(ctx, "/applications", createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Application](resp)
}

// UpdateApplication updates an application.
func (c *Client) UpdateApplication(ctx context.Context, name string, updateRequest *ApplicationUpdateRequest) (*Application, error) {
	resp, err := c.PutJSON(ctx, fmt.Sprintf("/application/%s", name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[Application](resp)
}

// DeleteApplication deletes an application.
func (c *Client) DeleteApplication(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor(fmt.Sprintf("/application/%s", name)).String(), nil)
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
