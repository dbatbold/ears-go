# Simple web resouce monitor
`ears-go` is a lightweight tool for monitoring web resources. It performs hourly checks and prints notifications when changes are detected in any resources specified in the `ears.json` configuration file. To prevent redundant alerts, the tool limits notifications to once per day for each changed resource until the `ears.json` file is updated. Configuration changes take effect immediately without requiring an application restart.

## Building and running the tool
```sh
$ make
$ ./ears &  # runs in background
```
