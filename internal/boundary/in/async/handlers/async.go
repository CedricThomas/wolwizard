package handlers

import (
	"context"
	"encoding/json"
	"fmt"

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
		if err := controller.WakeUpPCAgent(ctx, cmd.OSName); err != nil {
			return fmt.Errorf("wake up PC agent: %w", err)
		}
		return nil
	}
}

// ShutdownHost creates a typed callback for the PCAgent controller
func ShutdownHost(controller controller.PCAgent) async.Callback {
	return func(ctx context.Context, _ string) error {
		if err := controller.ShutdownCurrentHost(ctx); err != nil {
			return fmt.Errorf("shutdown host: %w", err)
		}
		return nil
	}
}
