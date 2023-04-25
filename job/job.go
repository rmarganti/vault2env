package job

import (
	"fmt"
	"os"

	"github.com/rmarganti/vault2env/config"
	"github.com/rmarganti/vault2env/sources"
)

type job interface {
	Run() error
}

func New(cfg *config.ConfigFile, originUri, targetUri, preset string) (job, error) {
	if preset == "init" {
		return &initJob{}, nil
	}

	// -[ Determine source URIs ]----------------------

	presetOrigin, presetTarget := getPresetSources(cfg, preset)

	if isInputPiped() {
		originUri = "std://"
	} else if originUri == "" {
		originUri = presetOrigin
	}

	if isOutputPiped() {
		targetUri = "std://"
	} else if targetUri == "" {
		targetUri = presetTarget
	}

	if originUri == "" || targetUri == "" {
		return &helpJob{}, nil
	}

	// -[ Build origin source ]------------------------

	originSrc, err := sources.New(originUri)

	if err != nil {
		return nil, fmt.Errorf("Error initializing origin source: %w", err)
	}

	// -[ Build target source ]------------------------

	targetSrc, err := sources.New(targetUri)

	if err != nil {
		return nil, fmt.Errorf("Error initializing target source: %w", err)
	}

	return &syncJob{
		origin: originSrc,
		target: targetSrc,
	}, nil
}

func getPresetSources(cfg *config.ConfigFile, preset string) (string, string) {
	if preset == "" {
		return "", ""
	}

	presetCfg, ok := cfg.Presets[preset]

	if !ok {
		return "", ""
	}

	return presetCfg.Origin, presetCfg.Target
}

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

// Determine if our output is being piped. Examples where this occurs:
//
// ```sh
// vault2env | cat
// vault2env > output.env
// ```
func isOutputPiped() bool {
	stdout, err := os.Stdout.Stat()

	if err != nil {
		return false
	}

	return (stdout.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
