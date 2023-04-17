package main

import (
	"fmt"
	"os"

	"github.com/rmarganti/vault2env/config"
	"github.com/rmarganti/vault2env/fetcher"
)

func main() {
	config, err := config.Load()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = fetcher.VaultToEnv(config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Wrote Vault secrets to .env")
}
