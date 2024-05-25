package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	// Configurar a aplicação
	app := application{}

	// Obter as rotas da aplicação
	mux := app.routes()

	// Imprimir mensagem indicando que o servidor está iniciando
	log.Println("Iniciando o servidor na porta 8080...")

	// Iniciar o servidor
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
