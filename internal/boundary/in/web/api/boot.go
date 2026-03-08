package api

import "errors"

type OSName string

const (
	OSWindows OSName = "windows"
	OSLinux   OSName = "linux"
)

type BootRequest struct {
	Name OSName `json:"name"`
}

func (r BootRequest) Validate() error {
	if r.Name == "" {
		return errors.New("os name is required")
	}

	switch r.Name {
	case OSWindows, OSLinux:
		return nil
	}

	return errors.New("unsupported os")
}

type BootData struct {
	Name OSName `json:"name"`
}

type BootResponse struct {
	Status bool     `json:"status"`
	Data   BootData `json:"data,omitempty"`
	Error  string   `json:"error,omitempty"`
}
