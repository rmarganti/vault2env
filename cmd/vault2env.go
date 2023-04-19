package main

import (
	"fmt"
	"os"

	"github.com/rmarganti/vault2env/config"
	"github.com/rmarganti/vault2env/puller"
	"github.com/rmarganti/vault2env/pusher"
)

func main() {
	config, err := config.Load()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	err = run(command, config)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(command string, config *config.Config) error {
	switch command {
	case "pull":
		return puller.VaultToEnv(config)

	case "push":
		return pusher.EnvToVault(config)

	case "help":
	default:
		showHelp()
		return nil
	}

	return nil
}

func showHelp() {
	fmt.Println("Usage: vault2env <command>")
	fmt.Println("Commands:")
	fmt.Println("  pull: Pull secrets from Vault and write to .env")
	fmt.Println("  push: Push secrets from .env to Vault")
	fmt.Println("  help: Show this help")
}
