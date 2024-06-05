package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// routes configura e retorna o roteador HTTP com as rotas e middlewares definidos
func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// registrar middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(app.addIPToContext)
	mux.Use(app.Session.LoadAndSave)

	// registrar rotas
	mux.Get("/", app.Home)
	mux.Post("/login", app.Login)

	mux.Get("/user/profile", app.Profile)

	// ativos estáticos
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
