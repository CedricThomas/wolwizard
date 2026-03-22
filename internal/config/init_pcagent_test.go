package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_PcAgent(t *testing.T) {
	t.Run("when all required env vars are set", func(t *testing.T) {
		// Given REDIS_URL and JWT_SECRET are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("PORT", "8081")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("PORT")
		})

		// When we initialize the PcAgent config
		cfg, err := Init(PcAgent)

		// Then we expect no error and correct values
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "redis://localhost:6379", cfg.RedisURL)
		assert.Equal(t, "test-secret", cfg.PcAgentConfig.JWTSecret)
		assert.Equal(t, "8081", cfg.PcAgentConfig.Port)
		assert.Equal(t, 86400, cfg.PcAgentConfig.JWTExpirySeconds)
		assert.Equal(t, "@every 5s", cfg.PcAgentConfig.MetricsReportingSchedule)
		assert.Equal(t, 5, cfg.PcAgentConfig.LastMetricsKeyTTLSeconds)
	})

	t.Run("when no optional env vars are set and defaults are used", func(t *testing.T) {
		// Given only required env vars are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("JWT_SECRET", "test-secret")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("JWT_SECRET")
		})

		// When we initialize the PcAgent config
		cfg, err := Init(PcAgent)

		// Then we expect default values
		assert.NoError(t, err)
		assert.Equal(t, "8081", cfg.PcAgentConfig.Port)
		assert.Equal(t, "@every 5s", cfg.PcAgentConfig.MetricsReportingSchedule)
		assert.Equal(t, 5, cfg.PcAgentConfig.LastMetricsKeyTTLSeconds)
		assert.Equal(t, 86400, cfg.PcAgentConfig.JWTExpirySeconds)
	})

	t.Run("when REDIS_URL is not set", func(t *testing.T) {
		// Given REDIS_URL is not set
		os.Unsetenv("REDIS_URL")
		os.Setenv("JWT_SECRET", "test-secret")
		t.Cleanup(func() {
			os.Unsetenv("JWT_SECRET")
		})

		// When we initialize the PcAgent config
		_, err := Init(PcAgent)

		// Then we expect an error containing REDIS_URL
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "REDIS_URL")
	})

	t.Run("when JWT_SECRET is not set", func(t *testing.T) {
		// Given REDIS_URL is set but JWT_SECRET is not
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Unsetenv("JWT_SECRET")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
		})

		// When we initialize the PcAgent config
		_, err := Init(PcAgent)

		// Then we expect an error containing JWT_SECRET
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "JWT_SECRET")
	})

	t.Run("when METRICS_REPORTING_SCHEDULE is invalid", func(t *testing.T) {
		// Given REDIS_URL, JWT_SECRET, and invalid METRICS_REPORTING_SCHEDULE are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("METRICS_REPORTING_SCHEDULE", "@invalid")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("METRICS_REPORTING_SCHEDULE")
		})

		// When we initialize the PcAgent config
		cfg, err := Init(PcAgent)

		// Then we expect an error and nil config
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("when LAST_METRICS_KEY_TTL_SECONDS is invalid", func(t *testing.T) {
		// Given REDIS_URL, JWT_SECRET, and invalid LAST_METRICS_KEY_TTL_SECONDS are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("LAST_METRICS_KEY_TTL_SECONDS", "invalid")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("LAST_METRICS_KEY_TTL_SECONDS")
		})

		// When we initialize the PcAgent config
		_, err := Init(PcAgent)

		// Then we expect an error
		assert.Error(t, err)
	})
}
