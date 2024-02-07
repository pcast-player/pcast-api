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
		println(fmt.Sprintf("Error: config file '%s' not found", file))
		os.Exit(1)
	}

	var cfg Config
	err = toml.Unmarshal(cfgString, &cfg)
	if err != nil {
		println(fmt.Sprintf("Error: config file '%s' is not valid %s", file, err))
		os.Exit(1)
	}

	return &cfg
}

func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
