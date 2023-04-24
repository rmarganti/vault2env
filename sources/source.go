package sources

import (
	"fmt"
	"net/url"
	"strings"
)

type Source interface {
	ReadSecrets() (secretsMap, error)
	WriteSecrets(secretsMap) error
}

func New(sourceUri string) (Source, error) {
	parsedUri, err := url.Parse(sourceUri)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse source URI: %w", err)
	}

	if parsedUri.Scheme == "" {
		parsedUri.Scheme = "file"
	}

	switch parsedUri.Scheme {
	case "file":
		return newFileSource(parsedUri.Host + parsedUri.Path), nil

	case "vault":
		return newVaultSource(parsedUri.Host, strings.Trim(parsedUri.Path, "/"))

	case "std":
		return newStdSource(), nil

	}

	return nil, fmt.Errorf("Unknown source type: %s", parsedUri.Scheme)
}
