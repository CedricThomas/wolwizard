package subscriptions

import (
	"context"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/boundary/in/async"
	asyncapi "github.com/CedricThomas/console/internal/boundary/in/async/api"
	"github.com/CedricThomas/console/internal/boundary/in/async/handlers"
	"github.com/CedricThomas/console/internal/controller"
)

// RegisterRaspberryAgent registers async channel subscriptions for Raspberry agent
func RegisterRaspberryAgent(
	ctx context.Context,
	consumer async.Consumer,
	raspberryController controller.RaspberryAgent,
) ([]func() error, error) {
	var unsubscribes []func() error

	unsubscribe, err := consumer.Subscribe(ctx, asyncapi.BootChannel, handlers.WakeUpPCAgent(raspberryController))
	if err != nil {
		return nil, fmt.Errorf("subscribe to boot channel: %w", err)
	}
	unsubscribes = append(unsubscribes, unsubscribe)

	log.Printf("Registered %d Raspberry async subscriptions", len(unsubscribes))
	return unsubscribes, nil
}
