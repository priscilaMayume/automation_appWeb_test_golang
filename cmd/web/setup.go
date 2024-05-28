package main

import (
	"os"
	"testing"
)

var app application

// TestMain é o teste principal para execução de testes
func TestMain (m *testing.M) {
	// Executa todos os testes no pacote
	os.Exit(m.Run())
}