package sources

import (
	"fmt"
	"os"
)

type fileSource struct {
	path string
}

func newFileSource(path string) *fileSource {
	return &fileSource{path: path}
}

func (src *fileSource) ReadSecrets() (secretsMap, error) {
	fmt.Fprintf(os.Stderr, "Reading secrets from %s\n", src.path)

	envFile, err := os.Open(src.path)

	if err != nil {
		return nil, fmt.Errorf("Error opening file: %w", err)
	}

	defer envFile.Close()

	secrets, err := NewSecretsFromEnv(envFile)

	if err != nil {
		return nil, fmt.Errorf("Error parsing file: %w", err)
	}

	return secrets, nil
}

func (src *fileSource) WriteSecrets(secrets secretsMap) error {
	fmt.Fprintf(os.Stderr, "Writing secrets to %s\n", src.path)

	envFile, err := os.Create(src.path)

	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}

	defer envFile.Close()

	err = secrets.ToEnv(envFile)

	if err != nil {
		return fmt.Errorf("Error writing to file: %w", err)
	}

	return nil
}
