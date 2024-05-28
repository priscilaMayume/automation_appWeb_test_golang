package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// register middlewere
	mux.Use(middleware.Recoverer)
	mux.Use(app.addIPToContext)

	// register routs
	mux.Get("/", app.Home)
	mux.Post("/login", app.Login)

	// static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}