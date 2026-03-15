package api

import "errors"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if r.Username == "" {
		return errors.New("'username' is required")
	}

	if r.Password == "" {
		return errors.New("'password' is required")
	}

	return nil
}
