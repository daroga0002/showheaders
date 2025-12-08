package config

import (
	"flag"
	"fmt"
)

// Config holds the application configuration
type Config struct {
	Port     int
	Hostname string
}

// Parse parses command line arguments and returns the configuration
func Parse() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.Port, "port", 8080, "Port to listen on")
	flag.StringVar(&cfg.Hostname, "hostname", "", "Hostname to bind to (empty for all interfaces)")
	flag.Parse()

	return cfg
}

// Address returns the formatted address string for the server
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}
