package base

import (
	"context"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/boundary/out/wol"
	"github.com/CedricThomas/console/internal/config"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain/async"
)

type rpAgent struct {
	wolSender wol.Sender
	config    *config.Config
}

func NewRaspberryAgentController(wolSender wol.Sender, cfg *config.Config) controller.RaspberryAgent {
	return &rpAgent{
		wolSender: wolSender,
		config:    cfg,
	}
}

func (ra *rpAgent) ExecuteBootMessage(ctx context.Context, bootMessage async.BootMessage) error {
	err := ra.wolSender.SendMagicPacket(ctx, ra.config.ServerNetworkAddress, ra.config.ServerMACAddress)
	if err != nil {
		return fmt.Errorf("send magic packet: %w", err)
	}
	log.Printf("Sent a wake up order for: %s\n", bootMessage.OSName)
	return nil
}
