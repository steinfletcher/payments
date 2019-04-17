package test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest/x/db"
	_ "github.com/steinfletcher/payments/migrations"
)

var DSN string

func tearDown(db *sqlx.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	db.MustExec(`TRUNCATE TABLE payments`)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

var Recorder *apitest.Recorder

func init() {
	//os.Setenv("DB_ADDR",  "postgresql://localhost:15432/payments?user=postgres&password=postgres&sslmode=disable")
	DSN = os.Getenv("DB_ADDR")
	Recorder = apitest.NewTestRecorder()
	wrappedDriver := db.WrapWithRecorder("postgres", Recorder)
	sql.Register("wrappedPostgres", wrappedDriver)
}

func DBSetup(setup func(tx *sqlx.Tx)) *sqlx.DB {
	d, err := goose.OpenDBWithDriver("postgres", DSN)
	if err != nil {
		panic(err)
	}

	err = goose.Up(d, ".")
	if err != nil {
		panic(err)
	}

	sqlxDB := sqlx.NewDb(d, DSN)

	tearDown(sqlxDB)

	tx, err := sqlxDB.Beginx()
	if err != nil {
		panic(err)
	}

	setup(tx)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return sqlxDB
}

func DBConnect() *sqlx.DB {
	testDB, err := sqlx.Connect("wrappedPostgres", DSN)
	if err != nil {
		panic(err)
	}
	return testDB
}

func SkipIntegration(t *testing.T) {
	if os.Getenv("DB_ADDR") == "" {
		t.Skip()
	}
}
