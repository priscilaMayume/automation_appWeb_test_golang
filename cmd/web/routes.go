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

	// register routs
	mux.Get("/", app.Home)

	// static assets

	return mux

}