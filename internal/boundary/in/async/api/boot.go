package async

import (
	"github.com/CedricThomas/console/internal/domain"
)

const BootChannel = "boot"

// BootCommand represents a boot request that can be published via pubsub
type BootCommand struct {
	OSName domain.OSName `json:"os_name"`
}
