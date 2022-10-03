package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/wombyz/amazon-dashboard/internal/config"
	"github.com/wombyz/amazon-dashboard/internal/handlers"
	"github.com/wombyz/amazon-dashboard/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var session *scs.SessionManager
var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	//gob.Register(models.Reservation{})
	gob.Register(map[string]int{})

	// change to true when in production
	app.InProduction = false
	app.UseCache = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)

	return nil
}
