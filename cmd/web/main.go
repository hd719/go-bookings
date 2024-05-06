package main

import (
	"encoding/gob"
	"flag"
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

const portNumber = ":8081"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the DB :)")
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring mail server..."))
	defer close(app.MailChan)
	listenForMail()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// Here we initialize what we are going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbName := flag.String("dbname", "", "Database Name")
	dbUser := flag.String("dbuser", "", "Database Password")
	dbPort := flag.String("dbport", "5432", "Database Port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")
	dbhost := flag.String("dbhost", "localhost", "Database Host")

	flag.Parse()

	// Create an email channel
	mailChan := make(chan models.MailData)
	// Add the channel to our App Config
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

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
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=secret sslmode=%s", *dbhost, *dbPort, *dbName, *dbUser, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	// db, err := driver.ConnectMongo("host=localhost port=27107 dbname=go-bookings user=system password=secret")
	if err != nil {
		log.Fatal("Cannot connect to DB!")
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	// This related to the the extra code in the handlers file on line 37, DO NOT DELETE
	// repo := handlers.NewRepo(&app, db)
	// handlers.NewHandlers(repo)

	// Note: the db connection is not tied to a specific database (pointer to a driver)
	handlers.NewRepo(&app, db)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
