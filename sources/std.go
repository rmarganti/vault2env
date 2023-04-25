package sources

import (
	"fmt"
	"os"
)

type stdSource struct {
}

func newStdSource() *stdSource {
	return &stdSource{}
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
