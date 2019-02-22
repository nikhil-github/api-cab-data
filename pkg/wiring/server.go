package wiring

import (
	"context"
	"github.com/nikhil-github/api-cab-data/pkg/service"
	"github.com/muesli/cache2go"
	"go.uber.org/zap"
	"net/http"
	"fmt"
	"github.com/pkg/errors"
	"github.com/nikhil-github/api-cab-data/pkg/cache"
	"github.com/nikhil-github/api-cab-data/pkg/database"
	"github.com/nikhil-github/api-cab-data/pkg/handler"
)

func StartServer(ctx context.Context, appName string,logger *zap.Logger) error {

	dbConfig := DatabaseConfig{
		URL: "root:password@tcp(localhost:3306)/cabtrips?parseTime=true",
	}

	db, err := NewDatabase(dbConfig)
	if err != nil {
		logger.Fatal("Failed to get database connection ",zap.Error(err))
	}

	cacheSvc := cache.New(cache2go.Cache("Cab-Trips-Data"))
	dbSvc := database.NewQueryer(db,logger)
	svc := service.NewService(dbSvc,cacheSvc,cacheSvc)
	router := NewRouter(&handler.Params{Logger:logger,Svc:svc,Cache:cacheSvc})

	errs := make(chan  error)
	serveHTTP(3000,logger,router,errs)

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return nil
	}
	return nil
}

func serveHTTP(port int, logger *zap.Logger, h http.Handler, errs chan error) {
	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{Addr: addr, Handler: h}

	go func() {
		logger.Info("Listening for HTTP requests", zap.String("http.address", addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- errors.Wrapf(err, "error serving HTTP on address %s", addr)
		}
	}()
}