package wiring

import (
	"github.com/go-chi/chi"
	"net/http"
	"go.uber.org/zap"
	"github.com/nikhil-github/api-cab-data/pkg/handler"
)


type Params struct {
	logger *zap.Logger
	svc    handler.TripServicer
	cache  handler.Clearer
}

func NewRouter(params Params) *chi.Mux {
	r := chi.NewRouter()

	//r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	//r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	//
	//r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/trips/pickupdate/{pickupDate}/medallion/{ids}", handler.Trips(params.logger,params.svc))
	r.Delete("/trip/clearcache", handler.ClearCache(params.logger,params.cache))


	return r
}