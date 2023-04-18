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

	switch command {
	case "pull":
		pull(config)

	case "push":
		push(config)

	case "help":
	default:
		showHelp()
	}

}

func pull(config *config.Config) {
	err := puller.VaultToEnv(config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func push(config *config.Config) {
	err := pusher.EnvToVault(config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Usage: vault2env <command>")
	fmt.Println("Commands:")
	fmt.Println("  pull: Pull secrets from Vault and write to .env")
	fmt.Println("  push: Push secrets from .env to Vault")
	fmt.Println("  help: Show this help")
}
