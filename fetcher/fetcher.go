package fetcher

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/rmarganti/vault2env/config"
)

// VaultToEnv fetches secrets from Vault and writes them to a .env file.
func VaultToEnv(cfg *config.Config) error {
	secrets, err := fetchSecrets(cfg)

	if err != nil {
		return err
	}

	err = writeEnv(secrets)

	return err
}

// Fetch the secrets from Vault.
func fetchSecrets(cfg *config.Config) (*vault.KVSecret, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Timeout = time.Second * 10

	client, err := vault.NewClient(vaultConfig)

	if err != nil {
		return nil, fmt.Errorf("Error creating Vault client: %w", err)
	}

	token, err := findToken()

	if err != nil {
		return nil, fmt.Errorf("Error finding Vault token: %w", err)
	}

	client.SetToken(token)
	secrets, err := client.KVv1(cfg.MountPath).Get(context.Background(), cfg.SecretPath)

	if err != nil {
		return nil, fmt.Errorf("Error fetching secret: %w", err)
	}

	return secrets, nil
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

// Write the secrets from Vault to the .env file.
func writeEnv(secrets *vault.KVSecret) error {
	envFile, err := os.Create(".env")

	if err != nil {
		return fmt.Errorf("Error opening .env file: %w", err)
	}

	sortedKeys := make([]string, 0, len(secrets.Data))
	for key := range secrets.Data {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		value := fmt.Sprint(secrets.Data[key])

		_, err = fmt.Fprintf(
			envFile,
			"%s=%s\n",
			key,
			strconv.Quote(value),
		)

		if err != nil {
			return fmt.Errorf("Error writing to .env file: %w", err)
		}
	}

	return nil
}
