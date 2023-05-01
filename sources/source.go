package sources

import (
	"fmt"
	"net/url"
)

type Source interface {
	// Read Secrets from the Source
	ReadSecrets() (secretsMap, error)

	// Write Secrets to the Source
	WriteSecrets(secretsMap) error
}

type SourceConstructor func(*url.URL) (Source, error)

// All supported Sources
var sources = map[string]SourceConstructor{
	"file":       newFileSourceFromURI,
	"k8s":        newK8sSourceFromURI,
	"kubernetes": newK8sSourceFromURI,
	"vault":      newVaultSourceFromURI,
	"std":        newStdSourceFromURI,
}

// New creates a new source from a URI. The URI is parsed,
// and the scheme is used to determine the source type.
func New(sourceUri string) (Source, error) {
	parsedUri, err := url.Parse(sourceUri)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse source URI: %w", err)
	}

	if parsedUri.Scheme == "" {
		parsedUri.Scheme = "file"
	}

	sourceFactory, ok := sources[parsedUri.Scheme]

	if !ok {
		return nil, fmt.Errorf("Unknown source type: %s", parsedUri.Scheme)
	}

	return sourceFactory(parsedUri)
}
