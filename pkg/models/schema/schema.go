package schema

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/wu-json/m1-finance-grafana/pkg/models/activity"
)

func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*activity.Activity)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
