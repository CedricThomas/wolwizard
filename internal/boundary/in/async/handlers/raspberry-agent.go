package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/boundary/in/async"
	"github.com/CedricThomas/console/internal/boundary/in/async/api"
	"github.com/CedricThomas/console/internal/controller"
)

// WakeUpPCAgent creates a typed callback for the RaspberryAgent controller
func WakeUpPCAgent(controller controller.RaspberryAgent) async.Callback {
	return func(ctx context.Context, rawCmd string) error {
		var cmd api.BootCommand
		if err := json.Unmarshal([]byte(rawCmd), &cmd); err != nil {
			return fmt.Errorf("invalid unmarshaling on consumption: %v", err)
		}
		log.Printf("Received boot command for OS: %s", cmd.OSName)
		if err := controller.WakeUpPCAgent(ctx, cmd.OSName); err != nil {
			return fmt.Errorf("wake up PC agent: %w", err)
		}
		return nil
	}
}
