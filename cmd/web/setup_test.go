package main

import (
	"os"
	"testing"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/repository/dbrepo"
)

// Variável global para a aplicação
var app application

// Função de teste principal
func TestMain(m *testing.M) {
	// Define o caminho para os templates
	pathToTemplates = "./../../templates/"
	
	// Inicializa a sessão da aplicação
	app.Session = getSession()
	app.DB = &dbrepo.TestDBRepo{}
	
	// Executa os testes
	os.Exit(m.Run())
}
