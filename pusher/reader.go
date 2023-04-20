package pusher

import (
	"fmt"
	"io"
	"os"

	envParse "github.com/hashicorp/go-envparse"
)

// Determine if our output is being piped. Examples where this occurs:
//
// ```sh
// cat .input.env | vault2env
// vault2env < input.env
// ````
func isInputPiped() bool {
	stdin, err := os.Stdin.Stat()

	if err != nil {
		return false
	}

	return (stdin.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}

// Read the secrets from the .env file.
func readSecretsFromFile() (*map[string]string, error) {
	fmt.Fprintln(os.Stderr, "Reading secrets from .env…")
	envFile, err := os.Open(".env")

	if err != nil {
		return nil, fmt.Errorf("Error opening .env file: %w", err)
	}

	return readSecrets(envFile)
}

// readSecrets reads lines of `key="value"` pairs
func readSecrets(r io.Reader) (*map[string]string, error) {
	secrets, err := envParse.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("Error parsing env: %w", err)
	}

	return &secrets, nil
}
