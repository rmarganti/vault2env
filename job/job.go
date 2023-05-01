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

	if isPiped(os.Stdin) {
		originUri = "std://"
	} else if originUri == "" {
		originUri = presetOrigin
	}

	if isPiped(os.Stdout) {
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

// Determine if a File is piped. Most likely, that File is actually
// a handler for Stdin or Stdout. Examples where this occurs:
//
// ```sh
// cat .input.env | vault2env
// vault2env < input.env
// ````
func isPiped(file *os.File) bool {
	fileInfo, err := file.Stat()

	if err != nil {
		return false
	}

	return (fileInfo.Mode() & os.ModeCharDevice) != os.ModeCharDevice
}
