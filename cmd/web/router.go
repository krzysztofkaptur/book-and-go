package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/krzysztofkaptur/book-and-go/internal/handlers"
)

func RunServer(repo *handlers.Repository) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(NoSurf)
	r.Use(SessionLoad)

	r.Get("/", repo.HomeHandler)
	r.Get("/about", repo.AboutHandler)
	r.Get("/search-availability", repo.AvailabilityHandler)
	r.Post("/search-availability", repo.PostAvailabilityHandler)
	r.Post("/search-availability-json", repo.AvailabilityJSONHandler)
	r.Get("/contact", repo.ContactHandler)
	r.Get("/generals-quarters", repo.GeneralsHandler)
	r.Get("/majors-suite", repo.MajorsHandler)
	r.Get("/make-reservation", repo.ReservationHandler)
	r.Post("/make-reservation", repo.PostReservationHandler)
	r.Get("/reservation-summary", repo.ReservationSummaryHandler)

	fileServer := http.FileServer(http.Dir("./static/"))

	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return r
}
