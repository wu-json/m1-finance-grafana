# m1-finance-grafana

Visualize your M1 Finance Account data using Grafana! Currently supported use-cases are:

- dividend visualization

## Inspiration

I personally love M1 Finance for dividend investing, but my one gripe with the platform is its lack of data visualizations for this use-case. However, since they allow you to download your dividend data as a CSV, I figured it wouldn't be difficult to build the visualizations I wanted myself using Grafana, which resulted in this project.

This project is super simple. It uses a Go script that parses and formats your M1 Finance dividend data CSVs and stores them in a local Postgres database (with TimescaleDB plugin), which can then be used as a data-source for a local Grafana instance.

## Requirements

- [Task v2.6.2](https://taskfile.dev/usage/) (`brew install go task`)
- [golang-migrate CLI v4.16.2](https://github.com/golang-migrate/migrate) (`brew install golang-migrate`)
- [sqlc v1.19.0](https://docs.sqlc.dev/en/stable/overview/install.html) (`brew install sqlc`)
- [Go v1.20.5](https://go.dev/doc/install)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

## Running Locally

First, you will want to prepare the database for ingesting your data. We will use Docker to spin up a Postgres instance and then apply the migrations using `golang-migrate`.

```bash
# spin up postgres & grafana locally
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

From there, you should see that the Postgres database contains your dividend data. You should now be able to visualize the data at the Grafana server hosted on `http://localhost:3000` (username is `user` and password is `pass` as configured in the `docker-compose.yaml`). The Postgres data-source and dividend dashboard should already be available as soon as you log in, as both of them are pre-provisioned.

## Generating Migrations

If you need to make updates to the database schema, you can create a new migration by running the following.

```bash
# generate a migration in ./database/migrations directory
task pg:migration:generate MIGRATION_NAME=my_migration
```
