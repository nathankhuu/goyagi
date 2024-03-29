package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/lob/logger-go"
	"github.com/nathankhuu/goyagi/pkg/application"
	"github.com/nathankhuu/goyagi/pkg/binder"
	"github.com/nathankhuu/goyagi/pkg/errors"
	"github.com/nathankhuu/goyagi/pkg/health"
	"github.com/nathankhuu/goyagi/pkg/metrics"
	"github.com/nathankhuu/goyagi/pkg/movies"
	"github.com/nathankhuu/goyagi/pkg/recovery"
	"github.com/nathankhuu/goyagi/pkg/signals"
)

// New returns a new HTTP server with the registered routes.
func New(app application.App) *http.Server {
	log := logger.New()

	e := echo.New()

	b := binder.New()
	e.Binder = b

	// Register our metrics middleware as our _first_ middlware
	e.Use(metrics.Middleware(app.Metrics))

	// Register the logger middleware after we set our custom binder
	e.Use(logger.Middleware())

	// Register our recovery middleware after our logger middleware
	e.Use(recovery.Middleware())

	errors.RegisterErrorHandler(e, app)

	health.RegisterRoutes(e)
	movies.RegisterRoutes(e, app)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
		Handler: e,
	}

	// signals.Setup() returns a channel we can wait until it's closed before we
	// shutdown our server
	graceful := signals.Setup()

	// start a goroutine that will wait for the graceful channel to close.
	// Becase this happens in a goroutine it will run concurrently with our
	// server but will not block the execution of this function.
	go func() {
		<-graceful
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Err(err).Error("server shutdown")
		}
	}()

	return srv
}
