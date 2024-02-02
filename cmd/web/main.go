package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Vikram222726/bookings/pkg/config"
	"github.com/Vikram222726/bookings/pkg/handlers"
	"github.com/Vikram222726/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

// main is the entry function of our application code
func main(){
	// should be set as true for production
	app.InProduction = false

	// by default session uses cookies to store our session data..
	// badgerstore is a built in db, which we can put in our go app and it will store all our sessions there..
	// other than this we can use memstore, mysqlstore, postgresql, redisstore and sqlite db's as well
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // should the cookies persist after the browser window is closed..
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict you want to be about which site this cookie is applied to..
	session.Cookie.Secure = app.InProduction // this is to ensure that the cookies are encrypted and connection is from https, must be true while running code in prod..
	
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.SetRepo(repo)

	render.NewTemplate(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	server := &http.Server{
		Addr: portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Server Started Listening on Port: %s...\n", portNumber)
	err = server.ListenAndServe()
	log.Fatal(err)

	// _ = http.ListenAndServe(portNumber, nil)
}