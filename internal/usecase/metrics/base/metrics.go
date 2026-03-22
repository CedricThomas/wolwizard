package base

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/CedricThomas/console/internal/domain"
	"github.com/CedricThomas/console/internal/service/keystore"
	ws "github.com/CedricThomas/console/internal/service/websocket"
	"github.com/CedricThomas/console/internal/usecase/metrics"
)

const (
	metricsKey = "metrics:last"
)

type metricsUsecase struct {
	lastMetricsKeyTTLSeconds time.Duration
	keystore                 keystore.Keystore
	wsManager                ws.Manager
}

func New(keystore keystore.Keystore, lastMetricsKeyTTLSeconds time.Duration, wsManager ws.Manager) metrics.Metrics {
	return &metricsUsecase{
		keystore:                 keystore,
		lastMetricsKeyTTLSeconds: lastMetricsKeyTTLSeconds,
		wsManager:                wsManager,
	}
}

func (m *metricsUsecase) ProcessMetrics(ctx context.Context, metrics domain.Metrics) error {
	log.Printf("Received metrics: OS %s, CPU %.2f%%, Memory %.2f%%, VRAM %.2f%%\n",
		metrics.OS, metrics.CPUUsage, metrics.MemoryUsage, metrics.VRAMUsage)

	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshal metrics: %w", err)
	}
	if err := m.keystore.SetWithTTL(ctx, metricsKey, string(metricsJSON), m.lastMetricsKeyTTLSeconds); err != nil {
		log.Printf("Error storing metrics: %v", err)
	}

	if err := m.wsManager.Broadcast(metricsJSON); err != nil {
		return fmt.Errorf("broadcast metrics to clients: %v", err)
	}

	return nil
}

func (m *metricsUsecase) GetLastMetrics(ctx context.Context) (*domain.Metrics, error) {
	// Check if metrics exist in keystore
	exists, err := m.keystore.Exists(ctx, metricsKey)
	if err != nil {
		return nil, fmt.Errorf("check metrics existence: %w", err)
	}
	if !exists {
		// If no metrics found, return nil without error
		return nil, nil
	}

	// Retrieve last metrics from keystore
	metricsJSON, err := m.keystore.Get(ctx, metricsKey)
	if err != nil {
		return nil, fmt.Errorf("get metrics: %w", err)
	}

	var data domain.Metrics
	if err := json.Unmarshal([]byte(metricsJSON), &data); err != nil {
		return nil, fmt.Errorf("unmarshal metrics: %w", err)
	}

	return &data, nil
}
