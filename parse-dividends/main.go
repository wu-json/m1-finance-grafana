package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/sqlc"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/utils"
)

func mapDividend(csvRecord []string) (sqlc.CreateDividendsParams, error) {
	if len(csvRecord) != 4 {
		return sqlc.CreateDividendsParams{}, fmt.Errorf("invalid csv record format: does not have 4 columns")
	}

	receivedOn, _ := time.Parse("Jan 2, 2006", csvRecord[0])
	ticker := strings.Split(csvRecord[2], " ")[0]
	dollarValue := csvRecord[3]
	activityType := csvRecord[1]

	validActivityTypes := []string{"Dividend", "Dividend - Deduction"}
	if !utils.Contains(validActivityTypes, activityType) {
		return sqlc.CreateDividendsParams{}, fmt.Errorf("invalid activity type: %s", activityType)
	}

	if activityType == "Dividend" {
		activityType = "Received"
	} else {
		activityType = "Deducted"
	}

	return sqlc.CreateDividendsParams{
		Ticker:       ticker,
		ActivityType: activityType,
		DollarValue:  sql.NullString{String: dollarValue, Valid: true},
		ReceivedOn:   receivedOn,
	}, nil
}

func validateHeaders(csvHeaders []string) error {
	if len(csvHeaders) != 4 {
		return fmt.Errorf("invalid csv headers: has %d columns instead of 4", len(csvHeaders))
	} else if csvHeaders[0] != "Date" {
		return fmt.Errorf("invalid csv headers: header 0 is %s instead of \"Date\"", csvHeaders[0])
	} else if csvHeaders[1] != "Activity" {
		return fmt.Errorf("invalid csv headers: header 1 is %s instead of \"Activity\"", csvHeaders[1])
	} else if csvHeaders[2] != "Summary" {
		return fmt.Errorf("invalid csv headers: header 2 is %s instead of \"Summary\"", csvHeaders[2])
	} else if csvHeaders[3] != "Value" {
		return fmt.Errorf("invalid csv headers: header 3 is %s instead of \"Value\"", csvHeaders[3])
	} else {
		return nil
	}
}

func processFile(ctx context.Context, queries *sqlc.Queries, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	allRecords, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	// validate headers
	if len(allRecords) < 1 {
		return fmt.Errorf("no records found")
	}
	err = validateHeaders(allRecords[0])
	if err != nil {
		return err
	}

	// ignore first record since that is the csv header
	records := allRecords[1:]

	for _, r := range records {
		dividend, err := mapDividend(r)
		if err != nil {
			return err
		}
		err = queries.CreateDividends(ctx, dividend)
		if err != nil {
			return err
		}

	}

	return nil
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
	filenames, err := utils.GetFileNames("../dividend-data")
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, file := range filenames {
		wg.Add(1)
		func() {
			defer wg.Done()
			fmt.Printf("Reading file: %s\n", file)
			err = processFile(ctx, queries, fmt.Sprintf("../dividend-data/%s", file))
		}()
	}

	wg.Wait()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
