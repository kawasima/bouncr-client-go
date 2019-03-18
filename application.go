package bouncr

import (
	"log"
	"encoding/json"
	"fmt"
	"net/http"
)

// Application application information
type Application struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PassTo      string `json:"pass_to"`
	VirtualPath string `json:"virtual_path"`
	TopPage     string `json:"top_page"`
}

// ApplicationSearchParams parameters for search applications
type ApplicationSearchParams struct {
	Offset int
	Limit  int
}

// ApplicationCreateRequest request for creating an application
type ApplicationCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PassTo      string `json:"pass_to,omitempty"`
	VirtualPath string `json:"virtual_path,omitempty"`
	TopPage     string `json:"top_page,omitempty"`
}

// ApplicationUpdateRequest request for creating an application
type ApplicationUpdateRequest ApplicationCreateRequest

// FindApplication find a application
func (c *Client) FindApplication(name string) (*Application, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/application/%s", name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *Application
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err

}

// ListApplications find the applications
func (c *Client) ListApplications(param *ApplicationSearchParams) ([]*Application, error) {
	req, err := http.NewRequest("GET", c.urlFor("/applications").String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(Application)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateApplication create an application
func (c *Client) CreateApplication(createRequest *ApplicationCreateRequest) (*Application, error) {
	resp, err := c.PostJSON("/applications", createRequest)
	defer closeResponse(resp)
	log.Printf("%s", resp.Body)

	if err != nil {
		return nil, err
	}

	var data Application

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateApplication update an application
func (c *Client) UpdateApplication(name string, updateRequest *ApplicationUpdateRequest) (*Application, error) {
	resp, err := c.PutJSON(fmt.Sprintf("/application/%s", name), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Application

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
