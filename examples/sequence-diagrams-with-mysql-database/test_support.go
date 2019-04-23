package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

const dbAddr = "mysql:mysql@/apitest"

func DBSetup(setup func(db *sqlx.DB)) *sqlx.DB {
	db, err := sqlx.Connect("mysql", dbAddr)
	if err != nil {
		panic(err)
	}

	goose.SetDialect("mysql")
	errMigration := goose.Up(db.DB, "./migrations")
	if errMigration != nil {
		panic(errMigration)
	}

	setup(db)
	return db
}
