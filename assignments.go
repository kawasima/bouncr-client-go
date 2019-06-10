package bouncr

import (
	"encoding/json"
	"net/http"
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

type Assignment struct {
	Group Group `json:"group"`
	Role  Role  `json:"role"`
	Realm Realm `json:"realm"`
}

// FindAssignments find assignments
func (c *Client) FindAssignment(findRequest *AssignmentRequest) (*Assignment, error) {
	url := c.urlFor("/assignment")
	q := url.Query()
	q.Set("group", findRequest.Group.Name)
	q.Set("role", findRequest.Role.Name)
	q.Set("realm", findRequest.Realm.Name)
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

	var data *Assignment
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateAssignments create assignments
func (c *Client) CreateAssignments(createRequest *[]AssignmentRequest) (*[]AssignmentRequest, error) {
	resp, err := c.PostJSON("/assignments", createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	return createRequest, nil
}

func (c *Client) DeleteAssignments(deleteRequest *[]AssignmentRequest) error {
	resp, err := c.requestJSON("DELETE",
		"/assignments",
		deleteRequest)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
