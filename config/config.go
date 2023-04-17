package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	MountPath string `json:"mountPath"`
	SecretPath string `json:"secretPath"`
}

// Load creates a new Config object from a `.vault2env.json`
// file in the current working directory.
func Load() (*Config, error) {
	configFile, err := os.Open(".vault2env.json")

	if err != nil {
		return nil, fmt.Errorf("Unable to load .vault2env.json: %w", err)
	}

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse .vault2env.json: %w", err)
	}

	if config.SecretPath == "" {
		return nil, fmt.Errorf("No `path` specified in .vault2env.json")
	}

	return &config, nil
}
