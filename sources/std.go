package sources

import (
	"fmt"
	"net/url"
	"os"
)

type stdSource struct {
}

func newStdSourceFromURI(uri *url.URL) (Source, error) {
	return &stdSource{}, nil
}

func (src *stdSource) ReadSecrets() (secretsMap, error) {
	secrets, err := NewSecretsFromEnv(os.Stdin)

	if err != nil {
		return nil, fmt.Errorf("Error parsing from stdin: %w", err)
	}

	return secrets, nil
}

func (src *stdSource) WriteSecrets(secrets secretsMap) error {
	err := secrets.ToEnv(os.Stdout)

	if err != nil {
		return fmt.Errorf("Error writing to stdout: %w", err)
	}

	return nil
}
