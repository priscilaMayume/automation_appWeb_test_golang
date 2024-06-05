package main

import (
	"log"
	"os"
	"testing"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/db"
)

// Variável global para a aplicação
var app application

// Função de teste principal
func TestMain(m *testing.M) {
	// Define o caminho para os templates
	pathToTemplates = "./../../templates/"
	
	// Inicializa a sessão da aplicação
	app.Session = getSession()
	app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)

	}
	defer conn.Close()
	app.DB = db.PostgresConn{DB : conn}

	// Executa os testes
	os.Exit(m.Run())
}
