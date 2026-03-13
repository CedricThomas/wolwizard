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

// RegisterWeb registers async channel subscriptions for Web agent
func RegisterWeb(
	ctx context.Context,
	consumer async.Consumer,
	webController controller.Web,
) ([]func() error, error) {
	var unsubscribes []func() error
	unsubscribe, err := consumer.Subscribe(ctx, asyncapi.MetricsChannel, handlers.ReportMetrics(webController))
	if err != nil {
		return nil, fmt.Errorf("subscribe to metrics channel: %w", err)
	}
	unsubscribes = append(unsubscribes, unsubscribe)

	log.Printf("Registered %d Web async subscriptions", len(unsubscribes))
	return unsubscribes, nil
}
