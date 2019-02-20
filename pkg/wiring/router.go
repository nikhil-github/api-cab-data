package wiring

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"time"
	"net/http"
	"github.com/nikhil-github/api-cab-data/pkg/cabs"
)

func route(svc *cabs.Service) {
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// write a health service
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/trips/pickupdate/{pickupDate}/medallion/{ids}", cabs.NewHandler(svc))

	http.ListenAndServe(":3000", r)
}