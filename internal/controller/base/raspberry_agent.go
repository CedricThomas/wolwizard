package base

import (
	"context"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/config"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/wol"
)

type rpAgent struct {
	wolSender wol.Sender
	config    *config.RaspberryAgentConfig
}

func NewRaspberryAgentController(wolSender wol.Sender, cfg *config.RaspberryAgentConfig) controller.RaspberryAgent {
	return &rpAgent{
		wolSender: wolSender,
		config:    cfg,
	}
}

func (ra *rpAgent) WakeUpPCAgent(ctx context.Context, osName domain.OSName) error {
	err := ra.wolSender.SendMagicPacket(ctx, ra.config.ServerNetworkAddress, ra.config.ServerMACAddress)
	if err != nil {
		return fmt.Errorf("send magic packet: %w", err)
	}
	log.Printf("Sent a wake up order for: %s\n", osName)
	return nil
}
