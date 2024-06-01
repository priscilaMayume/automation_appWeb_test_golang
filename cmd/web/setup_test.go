package main

import (
	"os"
	"testing"
)

// Variável global para a aplicação
var app application

// Função de teste principal
func TestMain(m *testing.M) {
	// Define o caminho para os templates
	pathToTemplates = "./../../templates/"
	
	// Inicializa a sessão da aplicação
	app.Session = getSession()

	// Executa os testes
	os.Exit(m.Run())
}
