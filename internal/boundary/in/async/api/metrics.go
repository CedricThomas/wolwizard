package api

const (
	MetricsChannel = "metrics/pc-agent"
)

type MetricsCommand struct {
	OS          string  `json:"os"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Uptime      string  `json:"uptime"`
	Status      string  `json:"status"`
}
