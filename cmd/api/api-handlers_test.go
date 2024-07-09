package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_app_authenticate(t *testing.T) {
	// Definição dos testes que serão realizados, cada teste com nome, corpo da requisição e código de status esperado.
	var theTests = []struct{
		name string
		requestBody string
		expectedStatusCode int
	}{
		{"valid user", `{"email":"admin@example.com","password":"secret"}`, http.StatusOK},
		{"not json", `I'm not JSON`, http.StatusUnauthorized},
		{"empty json", `{}`, http.StatusUnauthorized},
		{"empty email", `{"email":""}`, http.StatusUnauthorized},
		{"empty password", `{"email":"admin@example.com"}`, http.StatusUnauthorized},
		{"invalid user", `{"email":"admin@someotherdomain.com","password":"secret"}`, http.StatusUnauthorized},
	}

	// Loop pelos testes definidos
	for _, e := range theTests {
		var reader io.Reader
		// Cria um leitor com o corpo da requisição do teste atual
		reader = strings.NewReader(e.requestBody)
		// Cria uma nova requisição HTTP POST para a rota "/auth"
		req, _ := http.NewRequest("POST", "/auth", reader)
		// Cria um ResponseRecorder para gravar a resposta
		rr := httptest.NewRecorder()
		// Define o handler para a função de autenticação
		handler := http.HandlerFunc(app.authenticate)

		// Chama o handler com a requisição e o ResponseRecorder
		handler.ServeHTTP(rr, req)

		// Verifica se o status code retornado é o esperado
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}
