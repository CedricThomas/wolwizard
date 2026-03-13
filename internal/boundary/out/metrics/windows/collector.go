package windows

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CedricThomas/console/internal/boundary/out/metrics"
)

// collector implements the MetricsCollector interface for Windows systems
type collector struct{}

// New creates a new Windows metrics collector
func New() metrics.MetricsCollector {
	return &collector{}
}

// Collect gathers system metrics from the Windows system
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
		OS:          "Windows",
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		DiskUsage:   diskUsage,
		Uptime:      uptime,
		Status:      "running",
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

// getDiskUsage gets the disk usage percentage for the C: drive
func (c *collector) getDiskUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "wmic", "logicaldisk", "where", "DeviceID='\\\\.\\\\C:'", "get", "size,freespace")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("wmic command: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("disk info not found")
	}

	line := strings.TrimSpace(lines[1])
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, fmt.Errorf("parse disk: %w", fmt.Errorf("invalid format"))
	}

	sizeStr := strings.ReplaceAll(fields[0], ",", "")
	freeStr := strings.ReplaceAll(fields[1], ",", "")

	size, err1 := strconv.ParseFloat(sizeStr, 64)
	free, err2 := strconv.ParseFloat(freeStr, 64)

	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("parse disk values: %w", err1)
	}

	used := size - free
	return (used / size) * 100, nil
}

// getUptime gets the system uptime in a human-readable format
func (c *collector) getUptime(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-Command", "[string](Get-Date).Subtract((Get-CimInstance Win32_OperatingSystem).LastBootUpTime)")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("powershell command: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

var _ = time.Now()
var _ = regexp.MustCompile
