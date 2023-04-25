package job

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/rmarganti/vault2env/config"
)

type initJob struct{}

//go:embed .vault2env.default.json
var defaultConfig []byte

func (j *initJob) Run() error {
	if _, err := os.Stat(config.DefaultFilePath); err == nil {
		return fmt.Errorf(
			"Config file already exists at %s",
			config.DefaultFilePath,
		)
	}

	configFile, err := os.Create(config.DefaultFilePath)

	if err != nil {
		return fmt.Errorf("Error creating config file: %w", err)
	}

	defer configFile.Close()

	_, err = configFile.Write(defaultConfig)
	fmt.Fprintf(os.Stderr, "Wrote config file to %s\n", config.DefaultFilePath)

	if err != nil {
		return fmt.Errorf("Error writing config file: %w", err)
	}

	return nil
}
