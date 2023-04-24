package sources

import (
	"context"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
)

type vaultSource struct {
	client     *vault.Client
	mountPath  string
	secretPath string
}

func newVaultSource(mountPath, secretPath string) (*vaultSource, error) {
	token, err := findToken()

	if err != nil {
		return nil, err
	}

	client, err := newVaultClient(token)

	if err != nil {
		return nil, err
	}

	return &vaultSource{client, mountPath, secretPath}, nil
}

func (src *vaultSource) ReadSecrets() (secretsMap, error) {
	fmt.Fprintln(os.Stderr, "Reading secrets from Vault…")

	vaultSecrets, err := src.client.KVv1(src.mountPath).
		Get(context.Background(), src.secretPath)

	if err != nil {
		return nil, fmt.Errorf("Error fetching secret: %w", err)
	}

	// Convert all interface{} values to strings
	secrets := make(secretsMap, len(vaultSecrets.Data))
	for k, v := range vaultSecrets.Data {
		secrets[k] = fmt.Sprint(v)
	}

	return secrets, nil
}

func (src *vaultSource) WriteSecrets(secrets secretsMap) error {
	fmt.Fprintln(os.Stderr, "Writing secrets to Vault…")

	// Convert the map[string]string to map[string]interface{}
	vaultSecrets := make(map[string]interface{}, len(secrets))
	for k, v := range secrets {
		vaultSecrets[k] = v
	}

	err := src.client.KVv1(src.mountPath).
		Put(context.Background(), src.secretPath, vaultSecrets)

	if err != nil {
		return fmt.Errorf("Unable to write secrets to Vault: %w", err)
	}

	return nil
}

// Create a new Vault client instance
func newVaultClient(token string) (*vault.Client, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Timeout = time.Second * 10

	client, err := vault.NewClient(vaultConfig)

	if err != nil {
		return nil, fmt.Errorf("Error creating Vault client: %w", err)
	}

	client.SetToken(token)

	return client, nil
}

// Find a Vault token to use. First, check the VAULT_TOKEN environment variable.
// If that's not set, check the ~/.vault-token file.
func findToken() (string, error) {
	// Try the VAULT_TOKEN environment variable

	token := os.Getenv("VAULT_TOKEN")

	if token != "" {
		return token, nil
	}

	// Try to read from ~/.vault-token

	home, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("Error getting user home directory: %w", err)
	}

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
