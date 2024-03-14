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
