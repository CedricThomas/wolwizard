package windows

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/metrics"
)

// collector implements the MetricsCollector interface for Windows systems
type collector struct{}

// New creates a new Windows metrics collector
func New() metrics.Collector {
	return &collector{}
}

// Collect gathers system metrics from the Windows system
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
		OS:          domain.Windows,
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		VRAMUsage:   vramUsage,
	}, nil
}

// getCPUUsage gets the current CPU usage percentage
func (c *collector) getCPUUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "wmic", "cpu", "get", "loadpercentage")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("wmic command: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		if lineNum == 2 {
			loadStr := strings.TrimSpace(scanner.Text())
			load, err := strconv.ParseFloat(loadStr, 64)
			if err != nil {
				return 0, fmt.Errorf("parse load: %w", err)
			}
			return load, nil
		}
	}

	return 0, fmt.Errorf("cpu info not found")
}

// getMemoryUsage gets the current memory usage percentage
func (c *collector) getMemoryUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "wmic", "OS", "get", "FreePhysicalMemory,TotalVisibleMemorySize")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("wmic command: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("memory info not found")
	}

	lines[1] = strings.TrimSpace(lines[1])
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return 0, fmt.Errorf("parse memory: %w", fmt.Errorf("invalid format"))
	}

	freeMB, err1 := strconv.ParseFloat(fields[0], 64)
	totalKB, err2 := strconv.ParseFloat(fields[1], 64)

	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("parse memory values: %w", err1)
	}

	totalMB := totalKB / 1024
	usedMB := totalMB - freeMB
	return (usedMB / totalMB) * 100, nil
}

// getVRAMUsage gets the current VRAM usage percentage
func (c *collector) getVRAMUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "wmic", "path", "Win32_VideoController", "get", "AdapterRAM,CurrentAvailableVideoMemory", "/format:csv")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("wmic command: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// Skip header line
	scanner.Scan()

	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Split(line, ",")
		if len(fields) < 2 {
			return 0, fmt.Errorf("invalid video controller data format")
		}

		// AdapterRAM is in bytes, CurrentAvailableVideoMemory is in bytes
		totalVRAMStr := strings.Trim(fields[0], "\"")
		availableVRAMStr := strings.Trim(fields[1], "\"")

		totalVRAM, err1 := strconv.ParseFloat(totalVRAMStr, 64)
		availableVRAM, err2 := strconv.ParseFloat(availableVRAMStr, 64)

		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("parse VRAM values: %w", err1)
		}

		if totalVRAM == 0 {
			return 0, fmt.Errorf("total VRAM is zero")
		}

		usedVRAM := totalVRAM - availableVRAM
		return (usedVRAM / totalVRAM) * 100, nil
	}

	return 0, fmt.Errorf("video controller info not found")
}
