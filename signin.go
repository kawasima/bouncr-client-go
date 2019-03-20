package bouncr

import (
	"encoding/json"
)

// SignInRequest request for signing in
type SignInRequest struct {
	Account  string `json:"account,omitempty"`
	Password string `json:"password,omitempty"`
}

// SignInResponse response for signing in
type SignInResponse struct {
	Token string `json:"token,omitempty"`
}

// FindApplication find a application
func (c *Client) SignIn(signInRequest *SignInRequest) (*SignInResponse, error) {
	resp, err := c.PostJSON("/sign_in", signInRequest)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	var data *SignInResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}
