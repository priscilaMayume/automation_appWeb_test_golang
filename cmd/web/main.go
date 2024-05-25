package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

type application struct{}

func main() {
	// Configurar a aplicação
	app := application{}

	// Verificar se a porta 8080 está ocupada
	if isPortOpen(8080) {
		log.Println("A porta 8080 está ocupada. Fechando a porta...")
		closePort(8080)
	}

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

// Função para verificar se a porta está aberta
func isPortOpen(port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false // A porta está fechada
	}
	defer conn.Close()
	return true // A porta está aberta
}

// Função para fechar a porta
func closePort(port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return // A porta já está fechada
	}
	defer conn.Close()
}
