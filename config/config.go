package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultFilePath = ".vault2env.json"

type ConfigFile struct {
	Presets map[string]presetConfig `json:"presets"`
}

type presetConfig struct {
	Origin string `json:"from"`
	Target string `json:"to"`
}

// Load creates a new Config object from a
// config file in the current working directory.
func Load(path string) (*ConfigFile, error) {
	configFile, err := os.Open(path)

	if err != nil {
		// If the config file doesn't exist and is the
		// default config path, return an empty config.
		if path == DefaultFilePath {
			return &ConfigFile{}, nil
		}

		return nil, fmt.Errorf("Unable to load config: %w", err)
	}

	var config ConfigFile
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse %s: %w", path, err)
	}

	err = config.Validate()

	if err != nil {
		return nil, fmt.Errorf("Invalid config: %w", err)
	}

	return &config, nil
}

// TODO: Validate the Config structure
func (cfg *ConfigFile) Validate() error {
	return nil
}
