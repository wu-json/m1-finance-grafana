package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/format"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/sqlc"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/utils"
)

func processFile(ctx context.Context, logger *zap.SugaredLogger, queries *sqlc.Queries, path string) error {
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
	err = format.ValidateHeaders(allRecords[0])
	if err != nil {
		return err
	}

	// ignore first record since that is the csv header
	records := allRecords[1:]

	for _, r := range records {
		dividend, err := format.MapDividend(r)
		if err != nil {
			logger.Warnf("failed to map dividend: %v", r)
			continue
		}
		err = queries.CreateDividends(ctx, dividend)
		if err != nil {
			logger.Warnf("failed save dividend in pg: %v", r)
			continue
		}

	}

	return nil
}

func run(ctx context.Context, logger *zap.SugaredLogger) error {
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
		go func(file string) {
			defer wg.Done()
			if !strings.HasSuffix(file, ".csv") {
				logger.Warnf("not reading non-csv file :%s\n", file)
				return
			}
			logger.Infof("reading file: %s\n", file)
			err = processFile(ctx, logger, queries, fmt.Sprintf("../dividend-data/%s", file))
			if err != nil {
				logger.Warnf("error while reading file: %s\n", file, err)
			}
		}(file)
	}

	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	logger := zapLogger.Sugar()

	if err := run(ctx, logger); err != nil {
		logger.Error(err)
	}
}
