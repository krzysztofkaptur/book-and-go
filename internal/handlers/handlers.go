package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/krzysztofkaptur/book-and-go/internal/config"
	"github.com/krzysztofkaptur/book-and-go/internal/forms"
	"github.com/krzysztofkaptur/book-and-go/internal/models"
	render "github.com/krzysztofkaptur/book-and-go/internal/render"
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

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	remoteIP := repo.App.SessionManager.GetString(r.Context(), "remote_ip")

	stringMap := map[string]string{
		"test":      "Hello, again",
		"remote_ip": remoteIP,
	}

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) AvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	w.Write([]byte(fmt.Sprintf("%v - %v", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) AvailabilityJSONHandler(w http.ResponseWriter, r *http.Request) {
	response := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(response, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *Repository) ContactHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) GeneralsHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) MajorsHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationHandler(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := map[string]any{}
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := map[string]any{
			"reservation": reservation,
		}

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	m.App.SessionManager.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummaryHandler(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.SessionManager.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("Cannot get item from session")
		m.App.SessionManager.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.SessionManager.Remove(r.Context(), "reservation")

	data := map[string]any{
		"reservation": reservation,
	}

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
