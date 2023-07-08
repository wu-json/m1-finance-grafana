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

// Returns slice of strings containing filenames within the
// specified directory path.
func getFileNames(dirPath string) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return make([]string, 0), err
	}
	defer dir.Close()

	fis, err := dir.Readdir(-1)
	if err != nil {
		return make([]string, 0), err
	}

	data := make([]string, len(fis))
	for i, fi := range fis {
		data[i] = fi.Name()
	}

	return data, nil
}

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
	filenames, err := getFileNames("../dividend-data")
	if err != nil {
		return err
	}

	for _, file := range filenames {
		fmt.Printf("Reading file: %s\n", file)
		f, err := os.Open(fmt.Sprintf("../dividend-data/%s", file))
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
