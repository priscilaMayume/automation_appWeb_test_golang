package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

var app application
func init() {
    // Define o caminho relativo para os templates
    relativePath := "./../../templates/"
    absPath, err := filepath.Abs(relativePath)
    if err != nil {
        log.Fatalf("Error getting absolute path: %v", err)
    }
    pathToTemplates = absPath
}


// TestMain é a função principal de teste que configura o ambiente antes de executar os testes
func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"
	
	app.Session = getSession()

	// Executa os testes e retorna o resultado
	os.Exit(m.Run())
}
