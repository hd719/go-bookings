package render

import (
	"net/http"
	"testing"

	"github.com/hd719/go-bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("Failed, flash value of 123 not found")
	}
}

func TestRenderTemplate(t *testing.T) {
	// Note: app and pathToTemplates are package level variables can be overwritten
	pathToTemplates = "../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, r, "non-existing.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template doesnt exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	// Put session data in the context + headers
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	// Put the ctx with the session data back into the request
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app) // Note: app and pathToTemplates are package level variables can be overwritten
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}