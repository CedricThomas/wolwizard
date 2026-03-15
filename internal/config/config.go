package config

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	RedisURL                 string           `env:"REDIS_URL,required"`
	JWTSecret                string           `env:"JWT_SECRET"`
	JWTExpirySeconds         int              `env:"JWT_EXPIRY_SECONDS" envDefault:"86400"`
	ServerMACAddressStr      string           `env:"SERVER_MAC_ADDRESS"`
	ServerMACAddress         net.HardwareAddr `env:"-"` // Parsed from SERVER_MAC_ADDRESS
	ServerNetworkAddressStr  string           `env:"SERVER_NETWORK_ADDRESS"`
	ServerNetworkAddress     *net.UDPAddr     `env:"-"` // Parsed from SERVER_NETWORK_ADDRESS
	Port                     string           `env:"PORT"`
	MetricsReportingSchedule string           `env:"METRICS_REPORTING_SCHEDULE" envDefault:"@every 5s"`
}

// New creates a new Config instance with values from the environment
func New() (*Config, error) {
	var cfg Config
	var err error
	if err = env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("process env vars: %w", err)
	}

	if cfg.ServerMACAddressStr != "" {
		// Validate server MAC address at config initialization
		if cfg.ServerMACAddress, err = net.ParseMAC(cfg.ServerMACAddressStr); err != nil {
			return nil, fmt.Errorf("invalid MAC address in config: %w", err)
		}
	}

	if cfg.ServerNetworkAddressStr != "" {
		// Parse network address from SERVER_NETWORK_ADDRESS
		if cfg.ServerNetworkAddress, err = net.ResolveUDPAddr("udp", cfg.ServerNetworkAddressStr); err != nil {
			return nil, fmt.Errorf("invalid network address in config: %w", err)
		}
	}
	return &cfg, nil
}
