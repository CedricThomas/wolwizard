package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_Web(t *testing.T) {
	t.Run("when all required env vars are set", func(t *testing.T) {
		// Given REDIS_URL and JWT_SECRET are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("PORT", "8080")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("PORT")
		})

		// When we initialize the Web config
		cfg, err := Init(Web)

		// Then we expect no error and correct values
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "redis://localhost:6379", cfg.RedisURL)
		assert.Equal(t, "test-secret", cfg.WebConfig.JWTSecret)
		assert.Equal(t, "8080", cfg.WebConfig.Port)
		assert.Equal(t, 86400, cfg.WebConfig.JWTExpirySeconds)
		assert.Equal(t, 5, cfg.WebConfig.LastMetricsKeyTTLSeconds)
	})

	t.Run("when no env vars are set and defaults are used", func(t *testing.T) {
		// Given only required env vars are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("JWT_SECRET", "test-secret")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("JWT_SECRET")
		})

		// When we initialize the Web config
		cfg, err := Init(Web)

		// Then we expect default values
		assert.NoError(t, err)
		assert.Equal(t, "8080", cfg.WebConfig.Port)
		assert.Equal(t, 86400, cfg.WebConfig.JWTExpirySeconds)
		assert.Equal(t, 5, cfg.WebConfig.LastMetricsKeyTTLSeconds)
	})

	t.Run("when REDIS_URL is not set", func(t *testing.T) {
		// Given REDIS_URL is not set
		os.Unsetenv("REDIS_URL")
		os.Setenv("JWT_SECRET", "test-secret")
		t.Cleanup(func() {
			os.Unsetenv("JWT_SECRET")
		})

		// When we initialize the Web config
		_, err := Init(Web)

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

		// When we initialize the Web config
		_, err := Init(Web)

		// Then we expect an error containing JWT_SECRET
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "JWT_SECRET")
	})

	t.Run("when JWT_EXPIRY_SECONDS is invalid", func(t *testing.T) {
		// Given REDIS_URL, JWT_SECRET and invalid JWT_EXPIRY_SECONDS are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("JWT_SECRET", "test-secret")
		os.Setenv("JWT_EXPIRY_SECONDS", "invalid")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("JWT_EXPIRY_SECONDS")
		})

		// When we initialize the Web config
		_, err := Init(Web)

		// Then we expect an error
		assert.Error(t, err)
	})
}
