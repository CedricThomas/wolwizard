package api

import "errors"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	if r.Username == "" {
		return errors.New("'username' is required")
	}
	if r.Password == "" {
		return errors.New("'password' is required")
	}
	return nil
}

type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
