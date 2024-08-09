package handlers

import (
	"net/http"

	"github.com/krzysztofkaptur/book-and-go/pkg/config"
	"github.com/krzysztofkaptur/book-and-go/pkg/models"
	render "github.com/krzysztofkaptur/book-and-go/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	repo.App.SessionManager.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := repo.App.SessionManager.GetString(r.Context(), "remote_ip")

	stringMap := map[string]string{
		"test":      "Hello, again",
		"remote_ip": remoteIP,
	}

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
