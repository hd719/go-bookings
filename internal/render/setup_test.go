package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/hd719/go-bookings/internal/config"
	"github.com/hd719/go-bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func TestMain(m *testing.M) {
	// Here we initialize what we are going to store in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	// Creating Info Logger
	// Print logs to the terminal (stdout)
	infoLog = log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	// Creating Error Logger
	// Print logs to the terminal (stdout)
	errorLog = log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	testApp.Session = session

	// Sets the variable to app to the global variable app in render.go line 16
	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header // NOTE: Only have to return a type with a nil value (as long as the type is satisfied)
	return h
}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)

	return length, nil
}

func (tw *myWriter) WriteHeader(statusCode int) {

}
