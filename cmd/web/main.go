package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hd719/go-bookings/internal/config"
	"github.com/hd719/go-bookings/internal/driver"
	"github.com/hd719/go-bookings/internal/handlers"
	"github.com/hd719/go-bookings/internal/helpers"
	"github.com/hd719/go-bookings/internal/models"
	"github.com/hd719/go-bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	log.Println("Connected to the DB :)")

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// Here we initialize what we are going to store in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// Creating Info Logger
	// Print logs to the terminal (stdout)
	infoLog = log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// Creating Error Logger
	// Print logs to the terminal (stdout)
	errorLog = log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Create our session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Add our session to the Application config (global state)
	app.Session = session

	// Connect to DB
	log.Println("Connecting to DB...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go-bookings user=system password=secret")
	// db, err := driver.ConnectMongo("host=localhost port=27107 dbname=go-bookings user=system password=secret") ex.
	if err != nil {
		log.Fatal("Cannot connect to DB!")
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	// This related to the the extra code in the handlers file on line 37, DO NOT DELETE
	// repo := handlers.NewRepo(&app, db)
	// handlers.NewHandlers(repo)

	// Note: the db connection is not tied to a specific database (pointer to a driver)
	handlers.NewRepo(&app, db)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
