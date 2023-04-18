package pusher

import (
	"context"
	"fmt"
	"os"
	"time"

	envParse "github.com/hashicorp/go-envparse"
	vault "github.com/hashicorp/vault/api"
	"github.com/rmarganti/vault2env/config"
)

// EnvToVault fetches secrets from a .env file and writes them to Vault.
func EnvToVault(cfg *config.Config) error {
	secrets, err := fetchSecrets()

	if err != nil {
		return err
	}

	err = writeVault(cfg, secrets)

	if err != nil {
		return err
	}

	fmt.Println("Wrote .env secrets to Vault")

	return nil
}

// Fetch the secrets from the .env file.
func fetchSecrets() (*map[string]string, error) {
	fmt.Println("Reading secrets from .env…")
	envFile, err := os.Open(".env")

	if err != nil {
		return nil, fmt.Errorf("Error opening .env file: %w", err)
	}

	secrets, err := envParse.Parse(envFile)

	if err != nil {
		return nil, fmt.Errorf("Error parsing .env file: %w", err)
	}

	return &secrets, nil
}

// Write the secrets from the .env file to Vault.
func writeVault(cfg *config.Config, secrets *map[string]string) error {
	fmt.Println("Writing secrets to Vault…")
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Timeout = time.Second * 10

	client, err := vault.NewClient(vaultConfig)

	if err != nil {
		return fmt.Errorf("Error creating Vault client: %w", err)
	}

	// Convert the map[string]string to map[string]interface{}
	vaultSecrets := make(map[string]interface{}, len(*secrets))
	for k, v := range *secrets {
		vaultSecrets[k] = v
	}

	client.SetToken(cfg.Token)
	err = client.KVv1(cfg.MountPath).Put(context.Background(), cfg.SecretPath, vaultSecrets)

	if err != nil {
		return fmt.Errorf("Unable to write secrets to Vault: %w", err)
	}

	return nil
}
