package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct{
	Session *scs.SessionManager
}

func main() {
	// config do app
	app := application{}

	// obter a sessão
	app.Session = getSession()

	// print mensagem ao startar o serviço
	log.Println("Starting server on port 8080...")

	// iniciar o servidor
	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}