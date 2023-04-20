package puller

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	vault "github.com/hashicorp/vault/api"
)

// Determine if our output is being piped. Examples where this occurs:
//
// ```sh
// vault2env | cat
// vault2env > output.env
// ````
func isOutputPiped() bool {
	stdout, err := os.Stdout.Stat()

	if err != nil {
		return false
	}

	return (stdout.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

// Write the secrets from Vault to the .env file.
func writeEnvToFile(secrets *vault.KVSecret) error {
	fmt.Fprintln(os.Stderr, "Writing secrets to .env…")
	envFile, err := os.Create(".env")

	if err != nil {
		return fmt.Errorf("Error opening .env file: %w", err)
	}

	err = writeEnv(envFile, secrets)

	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Wrote Vault secrets to .env")

	return nil
}

// writeEnv writes Vault secrets as lines of `key="value"` pairs.
func writeEnv(w io.Writer, secrets *vault.KVSecret) error {
	sortedKeys := make([]string, 0, len(secrets.Data))
	for key := range secrets.Data {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		value := fmt.Sprint(secrets.Data[key])
		_, err := fmt.Fprintf(w, "%s=%s\n", key, strconv.Quote(value))

		if err != nil {
			return fmt.Errorf("Error writing to .env file: %w", err)
		}
	}

	return nil
}
