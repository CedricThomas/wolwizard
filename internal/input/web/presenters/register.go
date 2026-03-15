package presenters

import "github.com/CedricThomas/console/internal/input/web/api"

func RegisterSuccess() api.RegisterResponse {
	return api.RegisterResponse{
		Status:  "success",
		Message: "Account created successfully",
	}
}

func RegisterError(err error) api.RegisterResponse {
	if err == nil {
		return api.RegisterResponse{
			Status: "error",
			Error:  "unexpected error",
		}
	}
	return api.RegisterResponse{
		Status: "error",
		Error:  err.Error(),
	}
}
