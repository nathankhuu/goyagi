package main

import (
	"os"

	"github.com/lob/logger-go"
	"github.com/nathankhuu/goyagi/pkg/config"
	"github.com/nathankhuu/goyagi/pkg/database"
	migrations "github.com/robinjoseph08/go-pg-migrations"
)

const directory = "cmd/migrations"

func main() {
	log := logger.New()

	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		log.Err(err).Fatal("failed to connect to database")
	}

	err = migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Err(err).Fatal("failed migration")
	}
}
