# Simple web resource monitor
`ears-go` is a lightweight tool for monitoring web resources. It performs hourly checks and prints a notification when changes are detected in any resources specified in the `ears.json` configuration file. To prevent redundant alerts, the tool limits notifications to once per day for each changed resource until the `ears.json` file is updated. Configuration changes don't require an restart.

## Building and running the tool
```sh
$ go build ears.go
$ ./ears &  # runs in background
```
