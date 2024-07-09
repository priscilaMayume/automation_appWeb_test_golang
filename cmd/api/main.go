package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/repository"
	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/repository/dbrepo"
)

const port = 8090

type application struct {
	DSN string
	DB repository.DatabaseRepo
	Domain string
	JWTSecret string
}

func main() {
	// Estrutura da aplicação contendo as configurações necessárias
	var app application
	
	// Definindo as variáveis de linha de comando para domínio, DSN e segredo JWT
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for application, e.g. company.com")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Posgtres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "signing secret")
	flag.Parse()

	// Conectando ao banco de dados
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Inicializando o repositório do banco de dados
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	// Iniciando a API na porta definida
	log.Printf("Starting api on port %d\n", port)

	// Servindo as rotas da aplicação
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
