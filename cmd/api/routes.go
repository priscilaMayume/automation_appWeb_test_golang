package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Função que define as rotas da aplicação
func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// Registrando middleware
	mux.Use(middleware.Recoverer)
	// mux.Use(app.enableCORS) // Middleware para CORS (comentado)

	// Rotas de autenticação - manipuladores de autenticação e atualização de token
	mux.Post("/auth", app.authenticate)
	mux.Post("/refresh-token", app.refresh)

	// test handler
	mux.Get("/test", func(w http.ResponseWriter, r *http.Request){
		var payload = struct {
			Message string `json:"message"`
		}{
			Message: "hello, world",
		}

		_ = app.writeJSON(w, http.StatusOK, payload)
	})
	
	// Rotas protegidas
	mux.Route("/users", func(mux chi.Router) {
		// Uso do middleware de autenticação (comente esta linha para ativá-lo)
		
		mux.Get("/", app.allUsers) // Rota para obter todos os usuários
		mux.Get("/{userID}", app.getUser) // Rota para obter um usuário específico
		mux.Delete("/{userID}", app.deleteUser) // Rota para deletar um usuário
		mux.Put("/", app.insertUser) // Rota para inserir um novo usuário
		mux.Patch("/", app.updateUser) // Rota para atualizar um usuário existente
	})

	return mux
}
