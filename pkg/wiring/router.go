package wiring

import (
	"github.com/dimiro1/health"
	"github.com/gorilla/mux"
	"github.com/nikhil-github/api-cab-data/pkg/handler"
	"go.uber.org/zap"
)

type Params struct {
	Health health.Handler
	Logger *zap.Logger
	Svc    handler.TripServicer
	Cache  handler.Clearer
}

func NewRouter(params *Params) *mux.Router {
	rtr := mux.NewRouter()
	rtr.Handle("/trips/medallion/{ids}", handler.Trips(params.Logger, params.Svc)).Methods("GET")
	rtr.Handle("/trip/cache/contents", handler.ClearCache(params.Logger, params.Cache)).Methods("DELETE")
	rtr.Handle("/health", params.Health).Methods("GET")
	return rtr
}
