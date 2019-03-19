package bouncr

import (
	"encoding/json"
)

type IdObject struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// AssignmentRequest request for creating an assignment
type AssignmentRequest struct {
	Group IdObject `json:"group"`
	Role  IdObject `json:"role"`
	Realm IdObject `json:"realm"`
}

// CreateAssignments create assignments
func (c *Client) CreateAssignments(createRequest *[]AssignmentRequest) (*[]AssignmentRequest, error) {
	resp, err := c.PostJSON("/assignments", createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var data []AssignmentRequest

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) DeleteAssignments(deleteRequest *[]AssignmentRequest) error {
	resp, err := c.requestJSON("DELETE",
		c.urlFor("/assignments").String(),
		deleteRequest)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
