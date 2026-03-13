package presenters

import (
	"github.com/CedricThomas/console/internal/boundary/in/async/api"
	"github.com/CedricThomas/console/internal/domain"
)

func MetricsCommandToDomain(metrics api.MetricsCommand) domain.Metrics {
	return domain.Metrics{
		OS:          domain.OSName(metrics.OS),
		CPUUsage:    metrics.CPUUsage,
		VRAMUsage:   metrics.VRAMUsage,
		MemoryUsage: metrics.MemoryUsage,
	}
}

func DomainToMetricsCommand(metrics domain.Metrics) api.MetricsCommand {
	return api.MetricsCommand{
		OS:          string(metrics.OS),
		CPUUsage:    metrics.CPUUsage,
		VRAMUsage:   metrics.VRAMUsage,
		MemoryUsage: metrics.MemoryUsage,
	}
}
