package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// application estrutura que contém o gerenciador de sessões
type application struct {
	Session *scs.SessionManager
}

func main() {
	// Configura a aplicação
	app := application{}

	// Obtém um gerenciador de sessões
	app.Session = getSession()

	// Imprime uma mensagem
	log.Println("Starting server on port 8080...")

	// Inicia o servidor
	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
