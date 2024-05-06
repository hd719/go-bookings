package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/hd719/go-bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	// Testing the type that is return from routes() func
	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.mux and type is %T", v))
	}
}
