package presenters

import (
	"github.com/CedricThomas/console/internal/boundary/in/web/api"
	"github.com/CedricThomas/console/internal/domain"
)

func BootSuccess(name domain.OSName) api.BootResponse {
	return api.BootResponse{
		Status: true,
		Data: api.BootData{
			OSName: api.OSName(name),
		},
	}
}

func BootError(err error) api.BootResponse {
	return api.BootResponse{
		Status: false,
		Error:  err.Error(),
	}
}
