# Example API for testing with the fuzzer

Example API for the cloud native fuzzer. This API has endpoints for creating and getting todo lists. Also has some
simple auth endpoints.

## How to use

Generate OpenAPI documentation:

```sh
make swag
```

Build the Docker image:

```sh
IMAGE=todo-api make image
```

## Authentication

OpenAPI security:

```http
https://swagger.io/docs/specification/2-0/authentication/
```

Goswag example:

```http
https://github.com/swaggo/swag/blob/master/example/celler/main.go
```

Example OAuth2 implementation:

```http
https://www.ory.sh/hydra/docs/5min-tutorial
```

