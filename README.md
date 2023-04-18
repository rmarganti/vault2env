# vault2env

This utility will fetch secrets from a Hashicorp Vault server and write them to
a file in a format that can be sourced by a shell or dotenv.

## Configuration

1. Save a `.vault2env.json` file to your project's root directory. It is a JSON
   file that looks like the following:

    ```json
    {
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
