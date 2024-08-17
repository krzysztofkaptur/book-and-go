package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/krzysztofkaptur/book-and-go/internal/handlers"
)

func TestRunServer(t *testing.T) {
	var repo handlers.Repository

	server := RunServer(&repo)

	switch v := server.(type) {
	case *chi.Mux:
	default:
		t.Errorf("type is not *chi.Mux, but is %T", v)
	}
}
