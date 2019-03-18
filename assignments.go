package bouncr

import (
	"encoding/json"
)

type IdObject struct {
	ID int `json:"id"`
}

// AssignmentCreateRequest request for creating an assignment
type AssignmentCreateRequest struct {
	Group IdObject `json:"group"`
	Role  IdObject `json:"role"`
	Realm IdObject `json:"realm"`
}

// CreateAssignments create assignments
func (c *Client) CreateAssignments(createRequest *[]AssignmentCreateRequest) (*[]AssignmentCreateRequest, error) {
	resp, err := c.PostJSON("/assignments", createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var data []AssignmentCreateRequest

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
