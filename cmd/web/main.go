package main

import (
	"encoding/gob"
	"flag" // Importando o pacote data usando caminho relativo	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"
	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/repository"
	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/repository/dbrepo"
)

// application estrutura que contém o gerenciador de sessões
type application struct {
	DSN string
	DB repository.DatabaseRepo
	Session *scs.SessionManager
}

func main() {

	gob.Register(data.User{})
	// Configura a aplicação
	app := application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)

	}
	defer conn.Close()
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}


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
