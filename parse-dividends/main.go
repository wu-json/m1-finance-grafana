package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/sqlc"
)

func run() error {
	ctx := context.Background()

	// open pg connection
	db, err := sql.Open("postgres", "user=user password=pass dbname=m1finance sslmode=disable")
	if err != nil {
		return err
	}

	queries := sqlc.New(db)

	// get names of each file in dividend-data directory
	dir, err := os.Open("../dividend-data")
	if err != nil {
		return err
	}

	fis, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	dir.Close()

	for _, fi := range fis {
		fmt.Printf("Reading file: %s\n", fi.Name())
		f, err := os.Open(fmt.Sprintf("../dividend-data/%s", fi.Name()))
		if err != nil {
			return err
		}

		csvReader := csv.NewReader(f)
		records, err := csvReader.ReadAll()
		if err != nil {
			return err
		}

		fmt.Print(records)

		f.Close()
	}

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
