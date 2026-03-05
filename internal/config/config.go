package config

type Config struct {
}

// New creates a new Config instance with default values
func New() *Config {
	var cfg Config
	return &cfg
}
