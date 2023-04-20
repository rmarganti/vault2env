package puller

import (
	"context"
	"fmt"
	"os"
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

	if isOutputPiped() {
		err = writeEnv(os.Stdout, secrets)
	} else {
		err = writeEnvToFile(secrets)
	}

	if err != nil {
		return err
	}

	return nil
}

// Fetch the secrets from Vault.
func fetchSecrets(cfg *config.Config) (*vault.KVSecret, error) {
	fmt.Fprintln(os.Stderr, "Fetching secrets from Vault…")
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
