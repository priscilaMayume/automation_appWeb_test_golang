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
	mux.Use(middleware.Recoverer) // middleware para recuperação de panics
	mux.Use(app.addIPToContext)   // middleware para adicionar IP ao contexto
	mux.Use(app.Session.LoadAndSave) // middleware para carga e salvamento de sessão

	// registrar rotas
	mux.Get("/", app.Home)      // rota GET para a página inicial
	mux.Post("/login", app.Login) // rota POST para o login do usuário

	// ativos estáticos
	fileServer := http.FileServer(http.Dir("./static/")) // servidor de arquivos estáticos
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // rota para servir arquivos estáticos

	return mux
}
