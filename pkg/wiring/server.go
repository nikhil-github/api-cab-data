package wiring

import (
	"context"
	"fmt"

	"github.com/nikhil-github/api-cab-data/pkg/cabs"
	"github.com/muesli/cache2go"
)

func StartServer(ctx context.Context, appName string)  {

	fmt.Println("starting server")

	dbConfig := DatabaseConfig{
		URL: "root:password@tcp(localhost:3306)/cabtrips?parseTime=true",
	}

	db, err := NewDatabase(dbConfig)
	fmt.Println("db err",err)
	cacheSvc := cabs.NewCache(cache2go.Cache("Cab-Trips-Data"))
	dbSvc := cabs.NewQueryer(db)
	svc := cabs.NewService(dbSvc,cacheSvc,cacheSvc)
	cabs.NewHandler(svc)

	route(svc)
}

