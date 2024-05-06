package handlers

import (
	"net/http"

	"github.com/hd719/go-bookings/pkg/config"
	"github.com/hd719/go-bookings/pkg/models"
	"github.com/hd719/go-bookings/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	config := &Repository{
		App: a,
	}

	return config
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Remote ip address
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	templateData := &models.TemplateData{
		StringMap: map[string]string{
			"test":      "Hello, this is a test",
			"remote_ip": remoteIP,
		},
	}

	render.RenderTemplate(w, "about.page.tmpl", templateData)
}
