package main

import (
	"os"
	"testing"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/repository/dbrepo"
)

// Variável global para armazenar a aplicação.
var app application
var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE3MjAyMDc0MTksImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.CG7PENn1ieUbFCdMdXHHZ2V-djULER8oW9Q6bnZI-bM"

// TestMain é a função especial do pacote "testing" que é executada antes de qualquer teste.
func TestMain(m *testing.M) {
	// Configuração inicial da aplicação para os testes.
	app.DB = &dbrepo.TestDBRepo{} // Usando um repositório de banco de dados de teste.
	app.Domain = "example.com" // Definindo o domínio para a aplicação.
	app.JWTSecret = "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160" // Configurando o segredo JWT.

	// Executa todos os testes e finaliza com o código de saída do teste.
	os.Exit(m.Run())
}
