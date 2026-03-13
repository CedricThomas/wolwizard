package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/boundary/in/async"
	"github.com/CedricThomas/console/internal/controller"
)

// ShutdownHost creates a typed callback for the PCAgent controller
func ShutdownHost(controller controller.PCAgent) async.Callback {
	return func(ctx context.Context, _ string) error {
		log.Printf("Received shutdown command for current host")
		if err := controller.ShutdownCurrentHost(ctx); err != nil {
			return fmt.Errorf("shutdown host: %w", err)
		}
		return nil
	}
}
