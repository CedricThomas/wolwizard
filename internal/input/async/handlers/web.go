package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/async"
	"github.com/CedricThomas/console/internal/input/async/api"
	"github.com/CedricThomas/console/internal/input/async/presenters"
)

// ReportMetrics creates a typed callback for sending metrics to the Web controller
func ReportMetrics(controller controller.Web) async.Callback {
	return func(ctx context.Context, rawMetrics string) error {
		var metrics api.MetricsCommand
		if err := json.Unmarshal([]byte(rawMetrics), &metrics); err != nil {
			return fmt.Errorf("invalid unmarshaling of metrics: %v", err)
		}
		log.Printf("Received metrics: CPU %.2f%%, Memory %.2f%%", metrics.CPUUsage, metrics.MemoryUsage)

		if err := controller.ProcessMetrics(ctx, presenters.MetricsCommandToDomain(metrics)); err != nil {
			return fmt.Errorf("process metrics: %v", err)
		}
		return nil
	}
}
