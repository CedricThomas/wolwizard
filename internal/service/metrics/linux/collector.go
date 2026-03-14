package linux

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/metrics"
)

// collector implements the MetricsCollector interface for Linux systems
type collector struct{}

// New creates a new Linux metrics collector
func New() metrics.Collector {
	return &collector{}
}

// Collect gathers system metrics from the Linux system
func (c *collector) Collect(ctx context.Context) (domain.Metrics, error) {
	cpuUsage, err := c.getCPUUsage(ctx)
	if err != nil {
		return domain.Metrics{}, fmt.Errorf("get cpu usage: %w", err)
	}

	memoryUsage, err := c.getMemoryUsage(ctx)
	if err != nil {
		return domain.Metrics{}, fmt.Errorf("get memory usage: %w", err)
	}

	vramUsage, err := c.getVRAMUsage(ctx)
	if err != nil {
		return domain.Metrics{}, fmt.Errorf("get vram usage: %w", err)
	}

	return domain.Metrics{
		OS:          domain.Linux,
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		VRAMUsage:   vramUsage,
	}, nil
}

// getCPUUsage gets the current CPU usage percentage
func (c *collector) getCPUUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "top", "-bn1")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("top command: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	cpuRegex := regexp.MustCompile(`Cpu\s+([0-9.]+)\s*%us`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := cpuRegex.FindStringSubmatch(line); len(matches) > 1 {
			cpuUser, err := strconv.ParseFloat(matches[1], 64)
			if err != nil {
				return 0, fmt.Errorf("parse cpu user: %w", err)
			}

			idleRegex := regexp.MustCompile(`idle:?\s*([0-9.]+)`)
			if idleMatches := idleRegex.FindStringSubmatch(line); len(idleMatches) > 1 {
				idle, err := strconv.ParseFloat(idleMatches[1], 64)
				if err == nil {
					return 100 - idle, nil
				}
			}

			return cpuUser, nil
		}
	}

	return 0, fmt.Errorf("cpu usage not found")
}

// getMemoryUsage gets the current memory usage percentage
func (c *collector) getMemoryUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "free")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("free command: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		if lineNum == 2 {
			parts := strings.Fields(scanner.Text())
			if len(parts) >= 3 {
				total, err1 := strconv.ParseFloat(parts[1], 64)
				used, err2 := strconv.ParseFloat(parts[2], 64)

				if err1 != nil || err2 != nil {
					return 0, fmt.Errorf("parse memory values: %w", err1)
				}

				return (used / total) * 100, nil
			}
		}
	}

	return 0, fmt.Errorf("memory info not found")
}

// getVRAMUsage gets the current VRAM usage percentage from NVIDIA GPUs
func (c *collector) getVRAMUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "nvidia-smi", "--query-gpu=memory.used,memory.total", "--format=csv,nounits")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("nvidia-smi command: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("no GPU data found")
	}

	var totalUsed, totalMemory float64
	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], ",")
		if len(parts) < 2 {
			continue
		}

		used, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		total, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("parse VRAM values: %w", err1)
		}

		totalUsed += used
		totalMemory += total
	}

	if totalMemory == 0 {
		return 0, fmt.Errorf("total VRAM is zero")
	}

	return (totalUsed / totalMemory) * 100, nil
}
