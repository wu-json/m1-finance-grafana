package activity

import "time"

type Activity struct {
	ID           int64
	Activity     string
	Summary      string
	ValueDollars int64
	Date         time.Time
}
