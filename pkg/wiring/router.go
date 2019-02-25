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
	Svc    handler.Servicer
	Cache  handler.Clearer
}

// NewRouter configure all router.
func NewRouter(params *Params) *mux.Router {
	rtr := mux.NewRouter().StrictSlash(true)
	rtr.Handle("/trips/v1/medallions/{medallions}", handler.TripsByMedallion(params.Logger, params.Svc)).Methods("GET")
	rtr.Handle("/trips/v1/medallion/{medallion}/pickupdate/{pickupdate}", handler.TripsByMedallionsOnPickUpDate(params.Logger, params.Svc)).Methods("GET")
	rtr.Handle("/trips/v1/cache/contents", handler.ClearCache(params.Logger, params.Cache)).Methods("DELETE")
	rtr.Handle("/health", params.Health).Methods("GET")
	return rtr
}
