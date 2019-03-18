package bouncr

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Realm realm information
type Realm struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// RealmSearchParams parameters for search realms
type RealmSearchParams struct {
	Offset int
	Limit  int
}

// RealmCreateRequest request for creating an application
type RealmCreateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// RealmUpdateRequest request for creating an realm
type RealmUpdateRequest RealmCreateRequest

// FindRealm find a realm
func (c *Client) FindRealm(applicationName string, name string) (*Realm, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/application/%s/realm/%s", applicationName, name)).String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *Realm
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// ListApplications find the applications
func (c *Client) ListRealms(applicationName string, param *RealmSearchParams) ([]*Realm, error) {
	req, err := http.NewRequest("GET", c.urlFor(fmt.Sprintf("/application/%s/relms", applicationName)).String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(req)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data []*(Realm)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// CreateApplication create an application
func (c *Client) CreateRealm(applicationName string, createRequest *RealmCreateRequest) (*Realm, error) {
	resp, err := c.PostJSON(c.urlFor(fmt.Sprintf("/application/%s/relms", applicationName)).String(), createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var data Realm

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdateRealm update an realm
func (c *Client) UpdateRealm(applicationName string, updateRequest *RealmUpdateRequest) (*Realm, error) {
	resp, err := c.PutJSON(c.urlFor(fmt.Sprintf("/application/%s/relms", applicationName)).String(), updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data Realm
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
