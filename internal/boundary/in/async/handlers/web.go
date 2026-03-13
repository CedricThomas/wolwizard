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

// ReportMetrics creates a typed callback for sending metrics to the Web controller
func ReportMetrics(controller controller.Web) async.Callback {
	return func(ctx context.Context, rawMetrics string) error {
		var metrics api.MetricsCommand
		if err := json.Unmarshal([]byte(rawMetrics), &metrics); err != nil {
			return fmt.Errorf("invalid unmarshaling of metrics: %v", err)
		}
		log.Printf("Received metrics: CPU %.2f%%, Memory %.2f%%", metrics.CPUUsage, metrics.MemoryUsage)

		// TODO implement an api to domain converter and send a call to the Web controller to handle the received metrics
		return nil
	}
}
