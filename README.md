# m1-finance-grafana

This repository allows you to visualize M1 Finance Account data using Grafana.

## Requirements

- [Task](https://taskfile.dev/usage/)
- [Go v1.20.5](https://go.dev/doc/install)

## Running Locally

```bash
# spin up postgres locally
task docker:up

# configure workspace
go work init
go work use ./apps/parse-csv
```
