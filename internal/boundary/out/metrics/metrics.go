package metrics

import (
	"context"
)

// MetricsCollector defines the interface for collecting system metrics
type MetricsCollector interface {
	Collect(ctx context.Context) (*PCAgentMetrics, error)
}

// PCAgentMetrics represents the metrics collected from a PC agent
type PCAgentMetrics struct {
	OS          string  `json:"os"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Uptime      string  `json:"uptime"`
	Status      string  `json:"status"`
}
