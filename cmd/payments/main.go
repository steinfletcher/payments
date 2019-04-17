package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/steinfletcher/payments/api"
	"github.com/steinfletcher/payments/postgres"

	_ "github.com/lib/pq"
	_ "github.com/steinfletcher/payments/migrations"
)

type config struct {
	Port   string `env:"PORT" envDefault:"8080"`
	DBAddr string `env:"DB_ADDR"`
}

func main() {
	log.SetFlags(log.LstdFlags)

	// load config
	conf := &config{}
	err := env.Parse(conf)
	if err != nil {
		panic(err)
	}

	// perform database migrations
	db, err := goose.OpenDBWithDriver("postgres", conf.DBAddr)
	if err != nil {
		panic(err)
	}

	err = goose.Up(db, ".")
	if err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}

	// wire dependencies
	sqlxDB := sqlx.NewDb(db, conf.DBAddr)
	paymentsService := postgres.NewPaymentRepository(sqlxDB)

	// start server
	server := api.NewServer(paymentsService)
	log.Printf("Running server on :%s\n", conf.Port)
	server.Start(conf.Port)
}
