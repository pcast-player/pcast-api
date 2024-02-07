package config

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Host      string
	Port      int
	Logging   bool
	LogLevel  string `toml:"log_level"`
	LogFormat string `toml:"log_format"`
}

type Database struct {
	Logging bool
}

func New(file string) *Config {
	cfgString, err := os.ReadFile(file)
	if err != nil {
		panic("config.toml not found")
	}

	var cfg Config
	err = toml.Unmarshal(cfgString, &cfg)
	if err != nil {
		panic(fmt.Sprintf("error unmarshalling config.toml: %s", err))
	}

	return &cfg
}

func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
