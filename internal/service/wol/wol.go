package wol

//go:generate mockgen -source=wol.go -destination=mock/wol.go -package=mock -mock_names=Sender=MockSender
import (
	"context"
	"errors"
	"net"
)

var (
	ErrInvalidMAC = errors.New("invalid MAC address")
)

// WoL interface defines the operations for sending Wake-on-LAN magic packets.
type Sender interface {
	// SendMagicPacket sends a Wake-on-LAN magic packet to the specified MAC address on a specific network address.
	SendMagicPacket(ctx context.Context, networkAdddress *net.UDPAddr, mac net.HardwareAddr) error
}
