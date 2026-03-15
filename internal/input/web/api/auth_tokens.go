package api

import "errors"

type RevokeTokenRequest struct {
	Token string `json:"token"`
}

func (r RevokeTokenRequest) Validate() error {
	if r.Token == "" {
		return errors.New("'token' is required")
	}

	return nil
}
