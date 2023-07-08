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
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

## Running Locally

First, you will want to prepare the database for ingesting your data. We will use Docker to spin up a Postgres instance and then apply the migrations using `golang-migrate`.

```bash
# spin up postgres locally
task docker:up

# run migrations locally
task pg:migrate
```

From there, you will want to download `Dividend Activity` data from M1 Finance. You can do this by navigating to the account of interest on the website and then going to the `Activity Tab -> Filter Dividend Events -> Download CSV`. Make sure before you download, that you ensure the activity filter is set to filter for dividend events only.

Once you download all your CSVs, you will want to put them inside a directory inside the root project called `dividend-data`. I personally recommend you make each CSV contain 1 month of dividend data and then name the files something like `/m1-finance-grafana/dividend-data/2023-06-01-2023-07-01-dividends.csv`.

Now that we have the raw CSVs ready, we can ingest the data into Postgres by running the following.

```bash
# that's it!
task parse-dividends
```

From there, you should see that the Postgres database contains your dividend data.

## Generating Migrations

If you need to make updates to the database schema, you can create a new migration by running the following.

```bash
# generate a migration in ./database/migrations directory
task pg:migration:generate MIGRATION_NAME=my_migration
```
