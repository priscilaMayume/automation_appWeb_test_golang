package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// application estrutura que contém o gerenciador de sessões
type application struct {
	DSN string
	DB *sql.DB
	Session *scs.SessionManager
}

func main() {
	// Configura a aplicação
	app := application{}

	flag.StringVar(&app.DSN, "dns", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)


	}
	app.DB = conn


	// Obtém um gerenciador de sessões
	app.Session = getSession()

	// Imprime uma mensagem
	log.Println("Starting server on port 8080...")

	// Inicia o servidor
	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
