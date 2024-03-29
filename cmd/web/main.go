package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RorzGonzalez/gobook/pkg/config"
	"github.com/RorzGonzalez/gobook/pkg/handlers"
	"github.com/RorzGonzalez/gobook/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8009"

var app config.AppConfig

var session *scs.SessionManager

// main is the main application function
func main() {

	// change this to true when in InProduction
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	// Give render component access to our app templates
	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
