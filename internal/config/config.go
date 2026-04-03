package config

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
	"github.com/robfig/cron/v3"
)

type BinaryType int

const (
	Web BinaryType = iota
	PcAgent
	RaspberryAgent
)

type Config struct {
	RedisURL        string `env:"REDIS_URL,required"`
	WebConfig       WebConfig
	PcAgentConfig   PcAgentConfig
	RaspberryConfig RaspberryAgentConfig
}

type WebConfig struct {
	JWTSecret                string `env:"JWT_SECRET"`
	JWTExpirySeconds         int    `env:"JWT_EXPIRY_SECONDS" envDefault:"86400"`
	Port                     string `env:"PORT" envDefault:"8080"`
	LastMetricsKeyTTLSeconds int    `env:"LAST_METRICS_KEY_TTL_SECONDS" envDefault:"5"`
}

type PcAgentConfig struct {
	Port                     string `env:"PORT" envDefault:"8081"`
	JWTSecret                string `env:"JWT_SECRET"`
	JWTExpirySeconds         int    `env:"JWT_EXPIRY_SECONDS" envDefault:"86400"`
	MetricsReportingSchedule string `env:"METRICS_REPORTING_SCHEDULE" envDefault:"@every 5s"`
	LastMetricsKeyTTLSeconds int    `env:"LAST_METRICS_KEY_TTL_SECONDS" envDefault:"5"`
	BootOSTTLSeconds         int    `env:"BOOT_OS_TTL_SECONDS" envDefault:"300"`
}

type RaspberryAgentConfig struct {
	ServerMACAddress        net.HardwareAddr
	ServerMACAddressStr     string `env:"SERVER_MAC_ADDRESS"`
	ServerNetworkAddress    *net.UDPAddr
	ServerNetworkAddressStr string `env:"SERVER_NETWORK_ADDRESS"`
}

// Init creates a new Config instance and validates only the subconfig for the given binary type
func Init(bt BinaryType) (*Config, error) {
	var cfg Config
	var err error

	if err = env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("process env vars: %w", err)
	}

	switch bt {
	case Web:
		if err = validateWebConfig(&cfg.WebConfig); err != nil {
			return nil, fmt.Errorf("validate web config: %w", err)
		}
	case PcAgent:
		if err = validatePcAgentConfig(&cfg.PcAgentConfig); err != nil {
			return nil, fmt.Errorf("validate pc agent config: %w", err)
		}
	case RaspberryAgent:
		if err = validateRaspberryConfig(&cfg.RaspberryConfig); err != nil {
			return nil, fmt.Errorf("validate raspberry agent config: %w", err)
		}
	}

	return &cfg, nil
}

func validateWebConfig(cfg *WebConfig) error {
	if cfg.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func validatePcAgentConfig(cfg *PcAgentConfig) error {
	if cfg.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.MetricsReportingSchedule != "" {
		if _, err := cron.ParseStandard(cfg.MetricsReportingSchedule); err != nil {
			return fmt.Errorf("invalid METRICS_REPORTING_SCHEDULE: %w", err)
		}
	}
	if cfg.LastMetricsKeyTTLSeconds <= 0 {
		return fmt.Errorf("LAST_METRICS_KEY_TTL_SECONDS must be positive")
	}
	return nil
}

func validateRaspberryConfig(cfg *RaspberryAgentConfig) error {
	var err error

	if cfg.ServerMACAddressStr == "" {
		return fmt.Errorf("SERVER_MAC_ADDRESS is required")
	}

	if cfg.ServerMACAddress, err = net.ParseMAC(cfg.ServerMACAddressStr); err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}

	if cfg.ServerNetworkAddressStr == "" {
		return fmt.Errorf("SERVER_NETWORK_ADDRESS is required")
	}

	if cfg.ServerNetworkAddress, err = net.ResolveUDPAddr("udp", cfg.ServerNetworkAddressStr); err != nil {
		return fmt.Errorf("invalid network address: %w", err)
	}

	return nil
}
