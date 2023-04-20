# vault2env

Sync secrets between [Vault](https://www.vaultproject.io/) and a local `.env`.

## Installation

-   Binaries can be downloaded from the [releases page](https://github.com/rmarganti/vault2env/releases).
-   Mac users can install via Homebrew:

    ```sh
    brew tap rmarganti/tap
    brew install vault2env
    ```

## Configuration

1. Save a `.vault2env.json` file to your project's root directory. It is a JSON
   file that looks like the following:

    ```json
    {
        "$schema": "https://raw.githubusercontent.com/rmarganti/vault2env/main/vault2env.schema.json",
        "mountPath": "secret",
        "secretPath": "path/to/your/secrets"
    }
    ```

2. Ensure your Vault toke is stored either in the `VAULT_TOKEN` environment
   variable or in the `~/.vault-token` file.

## Pulling secrets from Vault

Your secrets will be downloaded and stored in the `.env` file in the current
working directory.

```sh
vault2env pull
```

## Pushing secrets to Vault

Secrets in your local `.env` will be stored in Vault.

```sh
vault2env push
```

## Options

| option     | description                                           |
| ---------- | ----------------------------------------------------- |
| `--config` | Path to a config file. Defaults to `.vault2env.json`. |
