package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hd719/go-bookings/internal/config"
	"github.com/hd719/go-bookings/internal/driver"
	"github.com/hd719/go-bookings/internal/forms"
	"github.com/hd719/go-bookings/internal/helpers"
	"github.com/hd719/go-bookings/internal/models"
	"github.com/hd719/go-bookings/internal/render"
	"github.com/hd719/go-bookings/internal/repository"
	"github.com/hd719/go-bookings/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) {
	config := &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
		// DB:  dbrepo.NewMongoRepo(db.Mongo, a),
	}

	// return config

	Repo = config
}

// This may not be needed, but DO NOT DELETE
// func NewHandlers(r *Repository) {
// 	Repo = r
// }

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Post Availability
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// Gets input values, by default the type is a string
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and the end is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	indent := "     "

	out, err := json.MarshalIndent(resp, "", indent)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), // Initializing an empty form when we go the reservation page
		Data: data,           // Create an empty reservation object when the page is first displayed
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get Reservation Data from the Post Req.
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	// PostForm contains the parsed form data from PATCH, POST or PUT body parameters.
	form := forms.New(r.PostForm)

	// Validation... (old)
	// Does the form have an value that is not an empty string
	// If the form has errors the Has func will create an error object
	// form.Has("first_name", r)

	// Validation... continued (new)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	// Form is not Valid:
	// Create the Form and Data fields that are going to be passed to TemplateData and get rendered on the client
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	// Form is Valid after passing validation:
	fmt.Println("The form is valid")

	// Adding the reservation object from line 121 into our session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// Http Redirect with a response code of 303
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Pulling out the reservation data from our session
	// Doing type assertion .(models.Reservation) forcefully asserting the type
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get Reservation from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Removing the reservation object from the our session storage because the reservation is now COMPLETE - we do not need it in our session storage in our client
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
