package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/krzysztofkaptur/book-and-go/internal/config"
	handlers "github.com/krzysztofkaptur/book-and-go/internal/handlers"
	"github.com/krzysztofkaptur/book-and-go/internal/models"
	"github.com/krzysztofkaptur/book-and-go/internal/render"
)

var app = config.AppConfig{}

func main() {
	repo, err := run()
	if err != nil {
		log.Fatal(err)
	}

	r := RunServer(repo)
	router := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = router.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*handlers.Repository, error) {
	gob.Register(models.Reservation{})

	// todo: change to true on production
	app.InProduction = false

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.SessionManager = session

	tc, err := render.CreateCacheTemplate()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	return repo, nil
}
