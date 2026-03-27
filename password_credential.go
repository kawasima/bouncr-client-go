package bouncr

import (
	"context"
	"net/http"
	"time"
)

// PasswordCredential represents a password credential.
type PasswordCredential struct {
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

// PasswordCredentialCreateRequest is a request for creating a password credential.
type PasswordCredentialCreateRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Initial  bool   `json:"initial"`
}

// PasswordCredentialUpdateRequest is a request for updating a password credential.
type PasswordCredentialUpdateRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// CreatePasswordCredential creates a password credential.
func (c *Client) CreatePasswordCredential(ctx context.Context, createRequest *PasswordCredentialCreateRequest) (*PasswordCredential, error) {
	resp, err := c.PostJSON(ctx, "/password_credential", createRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[PasswordCredential](resp)
}

// UpdatePasswordCredential updates a password credential.
func (c *Client) UpdatePasswordCredential(ctx context.Context, updateRequest *PasswordCredentialUpdateRequest) (*PasswordCredential, error) {
	resp, err := c.PutJSON(ctx, "/password_credential", updateRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	return decodeJSON[PasswordCredential](resp)
}

// DeletePasswordCredential deletes a password credential.
func (c *Client) DeletePasswordCredential(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.urlFor("/password_credential").String(), nil)
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
