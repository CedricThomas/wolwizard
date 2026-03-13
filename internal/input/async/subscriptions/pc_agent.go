package subscriptions

import (
	"context"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/async"
	asyncapi "github.com/CedricThomas/console/internal/input/async/api"
	"github.com/CedricThomas/console/internal/input/async/handlers"
)

// RegisterPCAgent registers async channel subscriptions for PC agent
func RegisterPCAgent(
	ctx context.Context,
	consumer async.Consumer,
	pcController controller.PCAgent,
) ([]func() error, error) {
	var unsubscribes []func() error

	// Subscribe to shutdown channel for PC agent
	if pcController != nil {
		unsubscribe, err := consumer.Subscribe(ctx, asyncapi.ShutdownChannel, handlers.ShutdownHost(pcController))
		if err != nil {
			return nil, fmt.Errorf("subscribe to shutdown channel: %w", err)
		}
		unsubscribes = append(unsubscribes, unsubscribe)
	}

	log.Printf("Registered %d PC async subscriptions", len(unsubscribes))
	return unsubscribes, nil
}
