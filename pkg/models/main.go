package main

import (
	"github.com/go-pg/pg/v10"
)

func main() {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	defer db.Close()
}
