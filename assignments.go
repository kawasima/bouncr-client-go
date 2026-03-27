package bouncr

import (
	"context"
	"net/http"
)

// IDObject represents an object identified by ID and name.
type IDObject struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// AssignmentRequest is a request for creating an assignment.
type AssignmentRequest struct {
	Group IDObject `json:"group"`
	Role  IDObject `json:"role"`
	Realm IDObject `json:"realm"`
}

// Assignment represents an assignment of a group, role, and realm.
type Assignment struct {
	Group Group `json:"group"`
	Role  Role  `json:"role"`
	Realm Realm `json:"realm"`
}

// FindAssignment finds an assignment.
func (c *Client) FindAssignment(ctx context.Context, findRequest *AssignmentRequest) (*Assignment, error) {
	u := c.urlFor("/assignment")
	q := u.Query()
	q.Set("group", findRequest.Group.Name)
	q.Set("role", findRequest.Role.Name)
	q.Set("realm", findRequest.Realm.Name)
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
	return decodeJSON[Assignment](resp)
}

// CreateAssignments creates assignments.
func (c *Client) CreateAssignments(ctx context.Context, createRequest *[]AssignmentRequest) (*[]AssignmentRequest, error) {
	resp, err := c.PostJSON(ctx, "/assignments", createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	return createRequest, nil
}

// DeleteAssignments deletes assignments.
func (c *Client) DeleteAssignments(ctx context.Context, deleteRequest *[]AssignmentRequest) error {
	resp, err := c.requestJSON(ctx, "DELETE", "/assignments", deleteRequest)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
