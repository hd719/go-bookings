package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	mux := routes()

	// Testing the type that is return from routes() func
	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.mux and type is %T", v))
	}
}
