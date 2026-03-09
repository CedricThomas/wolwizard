package api

import (
	"github.com/CedricThomas/console/internal/domain"
)

const (
	BootChannel     = "boot"
	ShutdownChannel = "shutdown"
)

// BootCommand represents a boot request that can be published via pubsub
type BootCommand struct {
	OSName domain.OSName `json:"os_name"`
}

// ShutdownCommand represents a shutdown request that can be published via pubsub
type ShutdownCommand struct{}

