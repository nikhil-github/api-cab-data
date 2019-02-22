package wiring

import (
	"github.com/nikhil-github/api-cab-data/pkg/handler"
	"github.com/gorilla/mux"
)



func NewRouter(params *handler.Params) *mux.Router {
	rtr := mux.NewRouter()
	rtr.Handle("/trips/pickupdate/{pickupDate}/medallion/{ids}", handler.Trips(params.Logger,params.Svc)).Methods("GET")
	rtr.Handle("/trip/cache/contents", handler.ClearCache(params.Logger,params.Cache)).Methods("DELETE")
	return rtr
}
