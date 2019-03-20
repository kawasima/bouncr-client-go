package bouncr

import (
	"encoding/json"
	"net/http"
	"time"
)

type PasswordCredential struct {
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
// PasswordCredentialCreateRequest request for creating an password credential
type PasswordCredentialCreateRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
// PasswordCredentialUpdateRequest request for updating an password credential
type PasswordCredentialUpdateRequest struct {
	Account     string `json:"account"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// CreatePasswordCredential create an password credential
func (c *Client) CreatePasswordCredential(createRequest *PasswordCredentialCreateRequest) (*PasswordCredential, error) {
	resp, err := c.PostJSON("/password_credential", createRequest)
	defer closeResponse(resp)

	if err != nil {
		return nil, err
	}

	var data PasswordCredential

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// UpdatePasswordCredential update an password credential
func (c *Client) UpdatePasswordCredential(updateRequest *PasswordCredentialUpdateRequest) (*PasswordCredential, error) {
	resp, err := c.PutJSON("/password_credential", updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data PasswordCredential
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeletePasswordCredential delete an password credential
func (c *Client) DeletePasswordCredential() error {
	req, err := http.NewRequest(
		"DELETE",
		c.urlFor("/password_credential").String(),
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
