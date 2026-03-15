package api

import "errors"

type OSName string

const (
	OSWindows OSName = "windows"
	OSLinux   OSName = "linux"
)

type BootRequest struct {
	OSName OSName `json:"os_name"`
}

func (r BootRequest) Validate() error {
	if r.OSName == "" {
		return errors.New("'os_name' is required")
	}

	switch r.OSName {
	case OSWindows, OSLinux:
		return nil
	}

	return errors.New("unsupported os")
}

type BootData struct {
	OSName OSName `json:"os_name"`
}

type BootResponse struct {
	Status bool     `json:"status"`
	Data   BootData `json:"data,omitempty"`
	Error  string   `json:"error,omitempty"`
}
