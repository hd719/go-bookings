package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// Do something (set vars) and then before exiting THEN RUN MY TESTS
	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
