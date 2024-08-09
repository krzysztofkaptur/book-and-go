package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/krzysztofkaptur/book-and-go/pkg/handlers"
)

func RunServer(repo *handlers.Repository) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(NoSurf)
	r.Use(SessionLoad)

	r.Get("/", repo.HomeHandler)
	r.Get("/about", repo.AboutHandler)

	router := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err := router.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
