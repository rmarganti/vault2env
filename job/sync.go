package job

import (
	"fmt"

	"github.com/rmarganti/vault2env/sources"
)

// A Job represents a single sync execution from an origin to a target.
type syncJob struct {
	origin sources.Source
	target sources.Source
}

// Run reads from the origin and writes to the target.
func (j *syncJob) Run() error {
	secrets, err := j.origin.ReadSecrets()

	if err != nil {
		return fmt.Errorf("error reading secrets: %w", err)
	}

	err = j.target.WriteSecrets(secrets)

	if err != nil {
		return fmt.Errorf("error writing secrets: %w", err)
	}

	return nil
}
