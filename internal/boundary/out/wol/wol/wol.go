package wol

import (
	"context"
	"fmt"
	"net"

	wolpkg "github.com/CedricThomas/console/internal/boundary/out/wol"
)

type wol struct{}

func New() wolpkg.Sender {
	return &wol{}
}

func (w *wol) SendMagicPacket(ctx context.Context, networkAdddress *net.UDPAddr, mac net.HardwareAddr) error {
	// Validate MAC address
	if len(mac) != 6 {
		return wolpkg.ErrInvalidMAC
	}

	// Create the magic packet
	magicPacket := make([]byte, 6+16*6)

	// 6 bytes of 0xFF
	for i := 0; i < 6; i++ {
		magicPacket[i] = 0xFF
	}

	// 16 repetitions of MAC
	for i := 0; i < 16; i++ {
		copy(magicPacket[6+i*6:], mac)
	}

	conn, err := net.DialUDP("udp", nil, networkAdddress)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(magicPacket)
	if err != nil {
		return fmt.Errorf("failed to send magic packet: %w", err)
	}

	return nil
}
