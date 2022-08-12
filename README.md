# HTTP Service

The goal of this project is to create a minimal and reliable HTTP service.

## Default Ports

### HTTP Port

The default transport protocol is HTTP and so no additional configuration is required to expose the service behind an HTTP server.

The service's defaut HTTP port is `8080`, and can be overridden by the `PORT` environment variable.

Ex. Use port `80` instead of `8080`.

```bash
export PORT='80'
```

## Logging

Each of the respective servers generates a log file. By default all logs are sent to STDOUT. This behvaviour and can be overridden by setting the `LOG_DIR` environment variable.

Ex. To store logs in `/var/log/go_server`.

```bash
export LOG_DIR='/var/log/accounts'
```

### Combined Logs

The service's default HTTP log file is `go_server.log`, and can be overridden by the `LOG_FILE` environment variable.

Ex. Log to file `logs.log` instead of `go_server.log`.

```bash
export LOG_FILE='logs.log'
```

## Run tests

```sh
go test ./... -coverprofile=cover.out && go tool cover -html=cover.out
```
