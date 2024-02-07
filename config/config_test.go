package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := New("./../fixtures/test/config.toml")
	if cfg == nil {
		t.Errorf("New() = nil, want *Config")
	}

	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, true, cfg.Server.Logging)
	assert.Equal(t, "log_level", cfg.Server.LogLevel)
	assert.Equal(t, "log_format", cfg.Server.LogFormat)
	assert.Equal(t, false, cfg.Database.Logging)
	assert.Equal(t, "localhost:3000", cfg.Server.GetAddress())
}
