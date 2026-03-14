package linux

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/metrics"
)

type collector struct{}

func New() metrics.Collector {
	return &collector{}
}

func (c *collector) Collect(ctx context.Context) (domain.Metrics, error) {
	cpuUsage, err := c.getCPUUsage(ctx)
	if err != nil {
		log.Printf("get cpu usage: %v", err)
	}

	memoryUsage, err := c.getMemoryUsage(ctx)
	if err != nil {
		log.Printf("get memory usage: %v", err)
	}

	vramUsage, err := c.getVRAMUsage(ctx)
	if err != nil {
		log.Printf("get vram usage: %v", err)
	}

	return domain.Metrics{
		OS:          domain.Linux,
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		VRAMUsage:   vramUsage,
	}, nil
}

func (c *collector) getCPUUsage(ctx context.Context) (float64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	line := strings.Split(string(data), "\n")[0]
	fields := strings.Fields(line)

	if len(fields) < 5 {
		return 0, fmt.Errorf("invalid /proc/stat format")
	}

	user, _ := strconv.ParseFloat(fields[1], 64)
	nice, _ := strconv.ParseFloat(fields[2], 64)
	system, _ := strconv.ParseFloat(fields[3], 64)
	idle, _ := strconv.ParseFloat(fields[4], 64)

	total := user + nice + system + idle
	if total == 0 {
		return 0, fmt.Errorf("invalid cpu total")
	}

	return (1 - idle/total) * 100, nil
}

func (c *collector) getMemoryUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx, "free", "-b")
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("free command: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "Mem:") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			return 0, fmt.Errorf("invalid memory format")
		}

		total, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("parse total memory: %w", err)
		}

		used, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return 0, fmt.Errorf("parse used memory: %w", err)
		}

		return (used / total) * 100, nil
	}

	return 0, fmt.Errorf("memory info not found")
}

func (c *collector) getVRAMUsage(ctx context.Context) (float64, error) {
	cmd := exec.CommandContext(ctx,
		"nvidia-smi",
		"--query-gpu=memory.used,memory.total",
		"--format=csv,noheader,nounits",
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("nvidia-smi command: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	var usedTotal float64
	var memoryTotal float64

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) != 2 {
			continue
		}

		used, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, fmt.Errorf("parse used vram: %w", err)
		}

		total, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return 0, fmt.Errorf("parse total vram: %w", err)
		}

		usedTotal += used
		memoryTotal += total
	}

	if memoryTotal == 0 {
		return 0, fmt.Errorf("total VRAM is zero")
	}

	return (usedTotal / memoryTotal) * 100, nil
}
