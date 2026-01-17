package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// DefaultJWTExpirationMin is the default JWT token expiration time in minutes
const DefaultJWTExpirationMin = 10

type Config struct {
	Server   Server
	Database Database
	Auth     Auth
}

type Auth struct {
	JwtSecret          string `toml:"jwt_secret"`
	JwtExpirationMin   int    `toml:"jwt_expiration_min"`
	GoogleClientID     string `toml:"google_client_id"`
	GoogleClientSecret string `toml:"google_client_secret"`
	GoogleRedirectURL  string `toml:"google_redirect_url"`
}

type Server struct {
	Host      string
	Port      int
	Logging   bool
	LogLevel  string `toml:"log_level"`
	LogFormat string `toml:"log_format"`
}

type Database struct {
	Host               string
	Port               int
	Database           string
	User               string
	Password           string
	MaxConnections     int    `toml:"max_connections"`
	MaxIdleConnections int    `toml:"max_idle_connections"`
	MaxLifetime        string `toml:"max_lifetime"`
	Logging            bool
	TimeZone           string `toml:"time_zone"`
}

func (d *Database) GetPostgresDSN() string {
	timeZone := d.TimeZone
	if timeZone == "" {
		timeZone = "UTC"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.Database, timeZone)
}

// New creates a new Config from a TOML file. Returns an error if the file cannot be read or parsed.
func New(file string) (*Config, error) {
	cfgString, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("config file '%s' not found: %w", file, err)
	}

	var cfg Config
	err = toml.Unmarshal(cfgString, &cfg)
	if err != nil {
		return nil, fmt.Errorf("config file '%s' is not valid: %w", file, err)
	}

	if cfg.Auth.JwtExpirationMin == 0 {
		cfg.Auth.JwtExpirationMin = DefaultJWTExpirationMin
	}

	return &cfg, nil
}

func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
