package async

import (
	"github.com/CedricThomas/console/internal/domain"
)

const BootChannel = "boot"

// BootMessage represents a boot request that can be published via pubsub
type BootMessage struct {
	OSName domain.OSName `json:"os_name"`
}
