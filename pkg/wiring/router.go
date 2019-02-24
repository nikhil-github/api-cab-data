package wiring

import (
	"github.com/dimiro1/health"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/nikhil-github/api-cab-data/pkg/handler"
)

type Params struct {
	Health health.Handler
	Logger *zap.Logger
	Svc    handler.TripServicer
	Cache  handler.Clearer
}

func NewRouter(params *Params) *mux.Router {
	rtr := mux.NewRouter()
	rtr.Handle("/trips/v1/medallion/{ids}", handler.Trips(params.Logger, params.Svc)).Methods("GET")
	rtr.Handle("/trips/v1/cache/contents", handler.ClearCache(params.Logger, params.Cache)).Methods("DELETE")
	rtr.Handle("/health", params.Health).Methods("GET")
	return rtr
}
