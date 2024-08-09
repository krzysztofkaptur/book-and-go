package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/krzysztofkaptur/book-and-go/pkg/config"
	handlers "github.com/krzysztofkaptur/book-and-go/pkg/handlers"
	"github.com/krzysztofkaptur/book-and-go/pkg/render"
)

var app = config.AppConfig{}

func main() {
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

	RunServer(repo)
}
