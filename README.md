# m1-finance-grafana

This repository allows you to visualize M1 Finance Account data using Grafana.

## Requirements

- [Task v2.6.2](https://taskfile.dev/usage/) (`brew install go task`)
- [golang-migrate CLI v4.16.2](https://github.com/golang-migrate/migrate) (`brew install golang-migrate`)
- [sqlc v1.19.0](https://docs.sqlc.dev/en/stable/overview/install.html) (`brew install sqlc`)
- [Go v1.20.5](https://go.dev/doc/install)

## Running Locally

```bash
# spin up postgres locally
task docker:up

# run migrations locally
task pg:migrate
```

## Generating Migrations

If you need to make updates to the database schema, you can create a new migration by running the following.

```bash
# generate a migration in ./database/migrations directory
task pg:migration:generate MIGRATION_NAME=my_migration
```
