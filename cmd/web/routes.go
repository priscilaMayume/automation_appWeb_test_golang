package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// Registrar middleware
	mux.Use(middleware.Recoverer)

	// Registrar rotas

	// Arquivos estáticos

	return mux
}
