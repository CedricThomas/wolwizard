package domain

type Metrics struct {
	OS          OSName  `json:"os"`
	CPUUsage    float64 `json:"cpu_usage"`
	VRAMUsage   float64 `json:"vram_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}
