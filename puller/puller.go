package puller

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

	if err != nil {
		return err
	}

	fmt.Println("Wrote Vault secrets to .env")

	return nil
}

// Fetch the secrets from Vault.
func fetchSecrets(cfg *config.Config) (*vault.KVSecret, error) {
	fmt.Println("Fetching secrets from Vault…")
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Timeout = time.Second * 10

	client, err := vault.NewClient(vaultConfig)

	if err != nil {
		return nil, fmt.Errorf("Error creating Vault client: %w", err)
	}

	client.SetToken(cfg.Token)
	secrets, err := client.KVv1(cfg.MountPath).Get(context.Background(), cfg.SecretPath)

	if err != nil {
		return nil, fmt.Errorf("Error fetching secret: %w", err)
	}

	return secrets, nil
}

// Write the secrets from Vault to the .env file.
func writeEnv(secrets *vault.KVSecret) error {
	fmt.Println("Writing secrets to .env…")
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
