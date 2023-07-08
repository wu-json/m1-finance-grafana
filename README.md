# m1-finance-grafana

Visualize your M1 Finance Account data using Grafana! Currently supported use-cases are:

- Dividend visualization

## Inspiration

I personally love M1 Finance for dividend investing, but my one gripe with the platform is its lack of data visualizations for this use-case. Since they allow you to download your dividend data as a CSV, I thought it would be nice to be able to create my own charts and tables using something like Grafana, which resulted in this project.

This project is super simple. It involves a Go script that parses and formats your M1 Finance dividend data CSVs and stores them in a local Postgres database, which can then be used as a data-source for a local Grafana instance.

> Please note that as a TypeScript one-trick, I'm still a Go novice. Do not
> judge my code too hard I promise it will get better...

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
