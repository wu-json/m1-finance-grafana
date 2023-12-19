# m1-finance-grafana

Visualize your M1 Finance Account data using Grafana! Currently supported use-cases are:

- dividend visualization

## Inspiration

I personally love M1 Finance for dividend investing, but my one gripe with the platform is its lack of data visualizations for this use-case. However, since they allow you to download your dividend data as a CSV, I figured it wouldn't be difficult to build the visualizations I wanted myself using Grafana, which resulted in this project.

This project is super simple. It uses a Go script that parses and formats your M1 Finance dividend data CSVs and stores them in a local Postgres database (with TimescaleDB plugin), which can then be used as a data-source for a local Grafana instance.

## Requirements

- [Task v2.6.2](https://taskfile.dev/usage/): `brew install go task`
- [golang-migrate CLI v4.16.2](https://github.com/golang-migrate/migrate): `brew install golang-migrate`
- [sqlc v1.19.0](https://docs.sqlc.dev/en/stable/overview/install.html): `brew install sqlc`
- [Go v1.20.5](https://go.dev/doc/install)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

## Running The Project

First, you will want to prepare the database for ingesting your data. We will use Docker to spin up a Postgres instance and then apply the migrations using `golang-migrate`.

```bash
# spin up postgres & grafana locally
task docker:up

# run migrations locally
task pg:migrate
```

From there, you will want to download `Dividend Activity` data from M1 Finance. You can do this by navigating to the account of interest on the website and then going to the `Activity Tab -> Filter Dividend Events -> Download CSV`. Make sure before you download, that you ensure the activity filter is set to filter for dividend events only.

Once you download all your CSVs, you will want to put them inside a directory inside the root project called `dividend-data`. I personally recommend you make each CSV contain 1 month of dividend data and then name the files something like `2023-12-page-1.csv` and back them up to cloud somewhere for easy access.

```bash
# directory structure
/m1-finance-grafana
    /dividend-data
        2023-05-01-2023-06-01-1.csv
        2023-05-01-2023-06-01-2.csv
        2023-06-01-2023-07-01-1.csv
    /parse-dividends
    ...
```

Now that we have the raw CSVs ready, we can ingest the data into Postgres by running the following.

```bash
# that's it!
task parse-dividends
```

From there, you should see that the Postgres database contains your dividend data. You should now be able to visualize the data at the Grafana server hosted on `http://localhost:3000` (username is `user` and password is `pass` as configured in the `docker-compose.yaml`). The Postgres data-source and dividend dashboard should already be available as soon as you log in, as both of them are pre-provisioned.

## Data Viz

Here is what the Dividends Dashboard looks like inside Grafana with some mock data. There are a couple of panels that display a time series of dividends earned, high earners, as well as a snapshot of the total dollar value of dividends earned during the period. This makes it much easier to predict dividend income over time, and can also help you figure out which dividend stocks to invest in to fill in any "dry" months in your dividend time series.

![Screenshot 2023-07-08 at 9 07 17 PM](https://github.com/wu-json/m1-finance-grafana/assets/45532884/e514fdda-9176-4e72-8454-bdedfb3195be)

## Generating Migrations

If you need to make updates to the database schema, you can create a new migration by running the following.

```bash
# generate a migration in ./database/migrations directory
task pg:migration:generate MIGRATION_NAME=my_migration
```
