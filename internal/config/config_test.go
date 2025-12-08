package config

import (
	"flag"
	"os"
	"testing"
)

func TestConfig_Address(t *testing.T) {
	testCases := []struct {
		name     string
		port     int
		hostname string
		expected string
	}{
		{
			name:     "default values",
			port:     8080,
			hostname: "",
			expected: ":8080",
		},
		{
			name:     "custom port",
			port:     3000,
			hostname: "",
			expected: ":3000",
		},
		{
			name:     "localhost",
			port:     8080,
			hostname: "localhost",
			expected: "localhost:8080",
		},
		{
			name:     "specific hostname and port",
			port:     9000,
			hostname: "192.168.1.1",
			expected: "192.168.1.1:9000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				Port:     tc.port,
				Hostname: tc.hostname,
			}

			result := cfg.Address()
			if result != tc.expected {
				t.Errorf("Address() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestParse_DefaultValues(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd"}

	cfg := Parse()

	if cfg.Port != 8080 {
		t.Errorf("default port = %v, want %v", cfg.Port, 8080)
	}

	if cfg.Hostname != "" {
		t.Errorf("default hostname = %v, want empty string", cfg.Hostname)
	}
}

func TestParse_CustomPort(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-port", "3000"}

	cfg := Parse()

	if cfg.Port != 3000 {
		t.Errorf("port = %v, want %v", cfg.Port, 3000)
	}
}

func TestParse_CustomHostname(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-hostname", "localhost"}

	cfg := Parse()

	if cfg.Hostname != "localhost" {
		t.Errorf("hostname = %v, want %v", cfg.Hostname, "localhost")
	}
}

func TestParse_AllOptions(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-port", "9000", "-hostname", "0.0.0.0"}

	cfg := Parse()

	if cfg.Port != 9000 {
		t.Errorf("port = %v, want %v", cfg.Port, 9000)
	}

	if cfg.Hostname != "0.0.0.0" {
		t.Errorf("hostname = %v, want %v", cfg.Hostname, "0.0.0.0")
	}
}
