package presenters

import "github.com/CedricThomas/console/internal/boundary/in/web/api"

func ShutdownSuccess() api.ShutdownResponse {
	return api.ShutdownResponse{
		Status: true,
	}
}

func ShutdownError(err error) api.ShutdownResponse {
	return api.ShutdownResponse{
		Status: false,
		Error:  err.Error(),
	}
}
