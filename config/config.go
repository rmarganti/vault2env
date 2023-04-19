package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	MountPath  string `json:"mountPath"`
	SecretPath string `json:"secretPath"`
	Token      string `json:"-"`
}

// Load creates a new Config object from a `.vault2env.json`
// file in the current working directory.
func Load(path *string) (*Config, error) {
	configFile, err := os.Open(*path)

	if err != nil {
		return nil, fmt.Errorf("Unable to load config: %w", err)
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

	token, err := findToken()

	if err != nil {
		return nil, fmt.Errorf("Error finding Vault token: %w", err)
	}

	config.Token = token

	return &config, nil
}

// Find a Vault token to use. First, check the VAULT_TOKEN environment variable.
// If that's not set, check the ~/.vault-token file.
func findToken() (string, error) {
	// Try the VAULT_TOKEN environment variable

	token := os.Getenv("VAULT_TOKEN")

	if token != "" {
		return token, nil
	}

	// Read the ~/.vault-token file
	home, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("Error getting user home directory: %w", err)
	}

	// Try to read from ~/.vault-token

	tokenFile, err := os.Open(home + "/.vault-token")

	if err != nil {
		return "", fmt.Errorf("Error opening Vault token file: %w", err)
	}

	_, err = fmt.Fscan(tokenFile, &token)

	if err != nil {
		return "", fmt.Errorf("Error reading Vault token file: %w", err)
	}

	return token, nil
}
