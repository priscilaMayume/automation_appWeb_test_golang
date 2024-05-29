package main

import (
	"os"
	"testing"
)

var app application // Variável global para o aplicativo

// TestMain é a função principal de teste que configura o ambiente antes de executar os testes
func TestMain(m *testing.M) {
	// Configura a sessão antes de executar os testes
	app.Session = getSession()

	// Executa os testes e retorna o resultado
	os.Exit(m.Run())
}
