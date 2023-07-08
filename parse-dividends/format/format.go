package format

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/sqlc"
	"github.com/wu-json/m1-finance-grafana/parse-dividends/utils"
)

// Transforms an M1 Finance CSV record into a format ready for Postgres insertion.
func MapDividend(csvRecord []string) (sqlc.CreateDividendsParams, error) {
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

// Ensures that the headers of an M1 Finance CSV are what we think they should be.
func ValidateHeaders(csvHeaders []string) error {
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
