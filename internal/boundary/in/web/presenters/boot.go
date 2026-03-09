package presenters

import (
	"github.com/CedricThomas/console/internal/boundary/in/web/api"
)

func BootSuccess(req api.BootRequest) api.BootResponse {
	return api.BootResponse{
		Status: true,
		Data: api.BootData{
			OSName: req.OSName,
		},
	}
}

func BootError(err error, req api.BootRequest) api.BootResponse {
	return api.BootResponse{
		Status: false,
		Data: api.BootData{
			OSName: req.OSName,
		},
		Error: err.Error(),
	}
}
