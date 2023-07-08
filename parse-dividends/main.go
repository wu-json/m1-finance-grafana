package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/sqlc"
)

func run() error {
	ctx := context.Background()

	db, err := sql.Open("postgres", "user=user password=pass dbname=m1finance sslmode=disable")
	if err != nil {
		return err
	}

	queries := sqlc.New(db)

	err = queries.CreateDividends(ctx, sqlc.CreateDividendsParams{
		Ticker:      "VOO",
		DollarValue: sql.NullString{String: "100.00", Valid: true},
		ReceivedOn:  time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
