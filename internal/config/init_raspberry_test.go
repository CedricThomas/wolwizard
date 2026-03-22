package config

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_RaspberryAgent(t *testing.T) {
	t.Run("when all required env vars are set", func(t *testing.T) {
		// Given REDIS_URL, SERVER_MAC_ADDRESS, and SERVER_NETWORK_ADDRESS are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("SERVER_MAC_ADDRESS", "30:56:0f:74:ff:2d")
		os.Setenv("SERVER_NETWORK_ADDRESS", "192.168.1.255:9")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("SERVER_MAC_ADDRESS")
			os.Unsetenv("SERVER_NETWORK_ADDRESS")
		})

		// When we initialize the RaspberryAgent config
		cfg, err := Init(RaspberryAgent)

		// Then we expect no error and correct values
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "redis://localhost:6379", cfg.RedisURL)
		assert.NotNil(t, cfg.RaspberryConfig.ServerMACAddress)
		assert.NotNil(t, cfg.RaspberryConfig.ServerNetworkAddress)
	})

	t.Run("when REDIS_URL is not set", func(t *testing.T) {
		// Given REDIS_URL is not set
		os.Unsetenv("REDIS_URL")
		os.Setenv("SERVER_MAC_ADDRESS", "30:56:0f:74:ff:2d")
		os.Setenv("SERVER_NETWORK_ADDRESS", "192.168.1.255:9")
		t.Cleanup(func() {
			os.Unsetenv("SERVER_MAC_ADDRESS")
			os.Unsetenv("SERVER_NETWORK_ADDRESS")
		})

		// When we initialize the RaspberryAgent config
		_, err := Init(RaspberryAgent)

		// Then we expect an error containing REDIS_URL
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "REDIS_URL")
	})

	t.Run("when SERVER_MAC_ADDRESS is not set", func(t *testing.T) {
		// Given REDIS_URL and SERVER_NETWORK_ADDRESS are set but SERVER_MAC_ADDRESS is not
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Unsetenv("SERVER_MAC_ADDRESS")
		os.Setenv("SERVER_NETWORK_ADDRESS", "192.168.1.255:9")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("SERVER_NETWORK_ADDRESS")
		})

		// When we initialize the RaspberryAgent config
		_, err := Init(RaspberryAgent)

		// Then we expect an error containing SERVER_MAC_ADDRESS
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SERVER_MAC_ADDRESS")
	})

	t.Run("when SERVER_NETWORK_ADDRESS is not set", func(t *testing.T) {
		// Given REDIS_URL and SERVER_MAC_ADDRESS are set but SERVER_NETWORK_ADDRESS is not
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("SERVER_MAC_ADDRESS", "30:56:0f:74:ff:2d")
		os.Unsetenv("SERVER_NETWORK_ADDRESS")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("SERVER_MAC_ADDRESS")
		})

		// When we initialize the RaspberryAgent config
		_, err := Init(RaspberryAgent)

		// Then we expect an error containing SERVER_NETWORK_ADDRESS
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SERVER_NETWORK_ADDRESS")
	})

	t.Run("when MAC address is invalid", func(t *testing.T) {
		// Given REDIS_URL, invalid SERVER_MAC_ADDRESS, and SERVER_NETWORK_ADDRESS are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("SERVER_MAC_ADDRESS", "invalid-mac")
		os.Setenv("SERVER_NETWORK_ADDRESS", "192.168.1.255:9")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("SERVER_MAC_ADDRESS")
			os.Unsetenv("SERVER_NETWORK_ADDRESS")
		})

		// When we initialize the RaspberryAgent config
		_, err := Init(RaspberryAgent)

		// Then we expect an error containing invalid MAC address
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid MAC address")
	})

	t.Run("when network address is invalid", func(t *testing.T) {
		// Given REDIS_URL, SERVER_MAC_ADDRESS, and invalid SERVER_NETWORK_ADDRESS are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("SERVER_MAC_ADDRESS", "30:56:0f:74:ff:2d")
		os.Setenv("SERVER_NETWORK_ADDRESS", "invalid-addr")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("SERVER_MAC_ADDRESS")
			os.Unsetenv("SERVER_NETWORK_ADDRESS")
		})

		// When we initialize the RaspberryAgent config
		_, err := Init(RaspberryAgent)

		// Then we expect an error containing invalid network address
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid network address")
	})

	t.Run("when MAC address is in different formats", func(t *testing.T) {
		// Given REDIS_URL and SERVER_NETWORK_ADDRESS are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("SERVER_NETWORK_ADDRESS", "192.168.1.255:9")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("SERVER_NETWORK_ADDRESS")
		})

		formats := []string{
			"30:56:0f:74:ff:2d",
			"30-56-0f-74-ff-2d",
			"3056.0f74.ff2d",
		}

		for _, mac := range formats {
			t.Run(mac, func(t *testing.T) {
				// Given SERVER_MAC_ADDRESS is set in a specific format
				os.Setenv("SERVER_MAC_ADDRESS", mac)
				t.Cleanup(func() {
					os.Unsetenv("SERVER_MAC_ADDRESS")
				})

				// When we initialize the RaspberryAgent config
				cfg, err := Init(RaspberryAgent)

				// Then we expect no error and correct values
				assert.NoError(t, err)
				assert.NotNil(t, cfg.RaspberryConfig.ServerMACAddress)
			})
		}
	})

	t.Run("when network address is broadcast or multicast", func(t *testing.T) {
		// Given REDIS_URL and SERVER_MAC_ADDRESS are set
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("SERVER_MAC_ADDRESS", "30:56:0f:74:ff:2d")
		t.Cleanup(func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("SERVER_MAC_ADDRESS")
		})

		addresses := []string{
			"192.168.1.255:9",   // Broadcast
			"224.0.0.251:9",     // Multicast
			"255.255.255.255:9", // Global broadcast
		}

		for _, addr := range addresses {
			t.Run(addr, func(t *testing.T) {
				// Given SERVER_NETWORK_ADDRESS is set to a broadcast or multicast address
				os.Setenv("SERVER_NETWORK_ADDRESS", addr)
				t.Cleanup(func() {
					os.Unsetenv("SERVER_NETWORK_ADDRESS")
				})

				// When we initialize the RaspberryAgent config
				cfg, err := Init(RaspberryAgent)

				// Then we expect no error and correct values
				assert.NoError(t, err)
				assert.NotNil(t, cfg.RaspberryConfig.ServerNetworkAddress)

				// Verify port is parsed correctly
				host, port, err := net.SplitHostPort(addr)
				assert.NoError(t, err)
				assert.Equal(t, "9", port)
				assert.Equal(t, host, cfg.RaspberryConfig.ServerNetworkAddress.IP.String())
				assert.Equal(t, 9, cfg.RaspberryConfig.ServerNetworkAddress.Port)
			})
		}
	})
}
