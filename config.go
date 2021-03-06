package main

import (
	"encoding/json"
	"os"
)

const (
	DefaultAddr    = "127.0.0.1"
	DefaultPort    = "8000"
	DefaultNatsURL = "0.0.0.0:4222"
)

// Config holds the necessary values needed to start the HTTP server
// connect to the NATS server and store the information.
type Config struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Certfile string `json:"certfile,omitempty"`
	Keyfile  string `json:"keyfile,omitempty"`
	NatsURL  string `json:"nats_url"`
}

// LoadConfig returns a Config instance from a given file name or
// an error if the file could not be read or decoded.
func LoadConfig(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}

	dec := json.NewDecoder(f)
	err = dec.Decode(cfg)

	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// DefaultConfig returns a Config instance with default values.
func DefaultConfig() *Config {
	return &Config{
		Address: DefaultAddr,
		Port:    DefaultPort,
		NatsURL: DefaultNatsURL,
	}
}
