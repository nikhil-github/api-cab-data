package wiring

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dimiro1/health"
	dbhealth "github.com/dimiro1/health/db"
	"github.com/muesli/cache2go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/nikhil-github/api-cab-data/pkg/cache"
	"github.com/nikhil-github/api-cab-data/pkg/database"
	"github.com/nikhil-github/api-cab-data/pkg/service"
)

func Start(cfg *Config, logger *zap.Logger) error {

	ctx := context.Background()
	dbx, err := NewDatabase(cfg.DB)
	if err != nil {
		logger.Fatal("Failed to get database connection ", zap.Error(err))
	}

	cacheSvc := cache.New(cache2go.Cache("Cab-Trips-Data"))
	dbSvc := database.NewQueryer(dbx, logger)
	svc := service.New(dbSvc, cacheSvc, cacheSvc, logger)
	router := NewRouter(&Params{Health: registerHealthCheck(dbx.DB), Logger: logger, Svc: svc, Cache: cacheSvc})

	errs := make(chan error)
	serveHTTP(cfg.HTTP.Port, logger, router, errs)

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return nil
	}
}

func serveHTTP(port int, logger *zap.Logger, h http.Handler, errs chan error) {
	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{Addr: addr, Handler: h}

	go func() {
		logger.Info("Listening for requests", zap.String("http.address", addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- errors.Wrapf(err, "error serving HTTP on address %s", addr)
		}
	}()
}

func registerHealthCheck(db *sql.DB) health.Handler {
	handler := health.NewHandler()
	handler.AddChecker("MySQL", dbhealth.NewMySQLChecker(db))
	return handler
}
