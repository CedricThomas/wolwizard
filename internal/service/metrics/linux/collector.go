package linux

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CedricThomas/console/internal/service/metrics"
)

// collector implements the MetricsCollector interface for Linux systems
type collector struct{}

// New creates a new Linux metrics collector
func New() metrics.MetricsCollector {
	return &collector{}
}

// Collect gathers system metrics from the Linux system
func (c *collector) Collect(ctx context.Context) (*metrics.PCAgentMetrics, error) {
	cpuUsage, err := c.getCPUUsage(ctx)
	if err != nil {
		return nil, fmt.Errorf("get cpu usage: %w", err)
	}

	memoryUsage, err := c.getMemoryUsage(ctx)
	if err != nil {
		return nil, fmt.Errorf("get memory usage: %w", err)
	}

	diskUsage, err := c.getDiskUsage(ctx)
	if err != nil {
		return nil, fmt.Errorf("get disk usage: %w", err)
	}

	uptime, err := c.getUptime(ctx)
	if err != nil {
		return nil, fmt.Errorf("get uptime: %w", err)
	}

	return &metrics.PCAgentMetrics{
		OS:          "Linux",
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		DiskUsage:   diskUsage,
		Uptime:      uptime,
		Status:      "running",
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

// getDiskUsage gets the disk usage percentage for the root partition
func (c *collector) getDiskUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "df", "-", "/")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("df command: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		if lineNum == 2 {
			parts := strings.Fields(scanner.Text())
			if len(parts) >= 5 {
				usageStr := strings.TrimSuffix(parts[4], "%")
				usage, err := strconv.ParseFloat(usageStr, 64)
				if err != nil {
					return 0, fmt.Errorf("parse disk usage: %w", err)
				}
				return usage, nil
			}
		}
	}

	return 0, fmt.Errorf("disk info not found")
}

// getUptime gets the system uptime in a human-readable format
func (c *collector) getUptime(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "uptime", "-p")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("uptime command: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

var _ = time.Now()
