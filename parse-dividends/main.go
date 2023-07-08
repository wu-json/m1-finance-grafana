package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
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
	defer db.Close()

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

		for i, r := range records {
			if i == 0 {
				continue
			}
			fmt.Println(r)

			receivedOn, _ := time.Parse("Jan 2, 2006", r[0])
			ticker := strings.Split(r[2], " ")[0]
			dollarValue := r[3]

			activityType := r[1]
			if activityType == "Dividend - Deduction" {
				activityType = "Deduction"
			} else if activityType != "Dividend" {
				return fmt.Errorf("invalid activity type: %s", activityType)
			}

			err = queries.CreateDividends(ctx, sqlc.CreateDividendsParams{
				Ticker:       ticker,
				ActivityType: activityType,
				DollarValue:  sql.NullString{String: dollarValue, Valid: true},
				ReceivedOn:   receivedOn,
			})
			if err != nil {
				return err
			}

		}

		fmt.Print(records)

		f.Close()
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
