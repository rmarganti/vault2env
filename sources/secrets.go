package sources

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	envParse "github.com/hashicorp/go-envparse"
)

type secretsMap map[string]string

// NewSecretsFromEnv create a secrets map from lines in the form of `KEY="VALUE"`.
func NewSecretsFromEnv(r io.Reader) (secretsMap, error) {
	secrets, err := envParse.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("Error parsing env: %w", err)
	}

	return secrets, nil
}

// NewSecretsFromInterfaceMap creates a secrets map from a map[string]interface{}.
func NewSecretsFromInterfaceMap(m map[string]interface{}) secretsMap {
	secrets := make(secretsMap, len(m))
	for k, v := range m {
		secrets[k] = fmt.Sprint(v)
	}
	return secrets
}

// NewSecretsFromByteMap creates a secrets map from a map[string][]byte.
func NewSecretsFromByteMap(m map[string][]byte) secretsMap {
	secrets := make(secretsMap, len(m))
	for k, v := range m {
		secrets[k] = string(v)
	}
	return secrets
}

// ToEnv writes secrets to a series of lines in the form `KEY=VALUE`.
func (s secretsMap) ToEnv(w io.Writer) error {
	sortedKeys := make([]string, 0, len(s))
	for key := range s {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		value := fmt.Sprint(s[key])
		_, err := fmt.Fprintf(w, "%s=%s\n", key, strconv.Quote(value))

		if err != nil {
			return fmt.Errorf("Error writing to .env file: %w", err)
		}
	}

	return nil
}

// ToInterfaceMap converts a secrets map to a map[string]interface{}.
// This is used, for example, when working with Vault.
func (s secretsMap) ToInterfaceMap() map[string]interface{} {
	secrets := make(map[string]interface{}, len(s))
	for k, v := range s {
		secrets[k] = v
	}
	return secrets
}
