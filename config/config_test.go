package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg, err := New("./../fixtures/test/config.toml")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, true, cfg.Server.Logging)
	assert.Equal(t, "log_level", cfg.Server.LogLevel)
	assert.Equal(t, "log_format", cfg.Server.LogFormat)
	assert.Equal(t, false, cfg.Database.Logging)
	assert.Equal(t, "localhost:3000", cfg.Server.GetAddress())
	assert.Equal(t, "host=localhost port=1337 user=pcast password=pcast dbname=pcast sslmode=disable TimeZone=Europe/Berlin", cfg.Database.GetPostgresDSN())
	assert.Equal(t, 1337, cfg.Database.Port)
	assert.Equal(t, "pcast", cfg.Database.Database)
	assert.Equal(t, "pcast", cfg.Database.User)
	assert.Equal(t, "pcast", cfg.Database.Password)
	assert.Equal(t, 10, cfg.Database.MaxConnections)
	assert.Equal(t, 5, cfg.Database.MaxIdleConnections)
	assert.Equal(t, "5m", cfg.Database.MaxLifetime)
}

func TestNew_FileNotFound(t *testing.T) {
	cfg, err := New("nonexistent.toml")
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "not found")
}

func TestNew_DefaultJWTExpiration(t *testing.T) {
	cfg, err := New("./../fixtures/test/config.toml")
	require.NoError(t, err)
	assert.Equal(t, DefaultJWTExpirationMin, cfg.Auth.JwtExpirationMin)
}
