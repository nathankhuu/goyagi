package application

import (
	"github.com/go-pg/pg"
	"github.com/nathankhuu/goyagi/pkg/config"
	"github.com/nathankhuu/goyagi/pkg/database"
	"github.com/nathankhuu/goyagi/pkg/metrics"
	"github.com/nathankhuu/goyagi/pkg/sentry"
	"github.com/pkg/errors"
)

// App contains necessary references that will be persisted throughout the
// application's lifecycle.
type App struct {
	Config  config.Config
	DB      *pg.DB
	Sentry  sentry.Sentry
	Metrics metrics.Metrics
}

// New creates a new instance of App with a Config and DB connection.
func New() (App, error) {
	cfg := config.New()

	db, err := database.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	sentry, err := sentry.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	m, err := metrics.New(cfg)
	if err != nil {
		return App{}, errors.Wrap(err, "application")
	}

	return App{cfg, db, sentry, m}, nil
}
