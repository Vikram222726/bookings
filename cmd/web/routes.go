package main

import (
	"net/http"

	"github.com/Vikram222726/bookings/pkg/config"
	"github.com/Vikram222726/bookings/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	// We'll be using chi module for routing as its lightweight, very fast
	// and provides us middleware in order to perform specific functions
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer) // Can use these middleware for Gracefully absorb panic and prints the stack trace
	mux.Use(WriteToConsole) // Custom middleware..
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	// Here we have to tell the server where to find the folder where our static files are stored
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}