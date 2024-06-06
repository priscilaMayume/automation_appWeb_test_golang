package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"
)

func Test_application_addIPToContext(t *testing.T) {
	tests := []struct{
		headerName string
		headerValue string
		addr string
		emptyAddr bool
	}{
		{"", "", "", false}, // Teste vazio
		{"", "", "", true}, // Teste com endereço vazio
		{"X-Forwarded-For", "192.3.2.1", "", false}, // Teste com cabeçalho X-Forwarded-For
		{"", "", "hello:world", false}, // Teste com endereço específico
	}

	// Cria um manipulador fictício que usaremos para verificar o contexto
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Certifica-se de que o valor existe no contexto
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Error(contextUserKey, "não está presente")
		}

		// Certifica-se de que recebemos uma string de volta
		ip, ok := val.(string)
		if !ok {
			t.Error("não é uma string")
		}
		t.Log(ip)
	})

	for _, e := range tests {
		// Cria o manipulador para teste
		handlerToTest := app.addIPToContext(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)

		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headerName) > 0 {
			req.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_application_ipFromContext(t *testing.T) {
	// Obtém um contexto
	ctx := context.Background()

	// Coloca algo no contexto
	ctx = context.WithValue(ctx, contextUserKey, "qualquer coisa")

	// Chama a função
	ip := app.ipFromContext(ctx)

	// Executa o teste
	if !strings.EqualFold("qualquer coisa", ip) {
		t.Error("valor incorreto retornado do contexto")
	}
}

func Test_app_auth(t *testing.T) {
	// Definindo um próximo manipulador de teste que não faz nada
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

	})

	// Definindo os casos de teste com dois cenários: "logged in" e "not logged in"
	var tests = []struct{
		name string
		isAuth bool
	}{
		{"logged in", true},
		{"not logged in", false},
	}

	// Looping através de cada caso de teste
	for _, e := range tests {
		// Obter o manipulador a ser testado
		handlerToTest := app.auth(nextHandler)
		// Criar uma nova requisição de teste
		req := httptest.NewRequest("GET", "http://testing", nil)
		// Adicionar contexto e sessão à requisição
		req = addContextAndSessionToRequest(req, app)
		// Se o caso de teste for autenticado, colocar o usuário na sessão
		if e.isAuth {
			app.Session.Put(req.Context(), "user", data.User{ID: 1})
		}
		// Criar um novo gravador de resposta de teste
		rr := httptest.NewRecorder()
		// Chamar o manipulador com a requisição e o gravador de resposta
		handlerToTest.ServeHTTP(rr, req)

		// Verificar se o código de status está correto quando autenticado
		if e.isAuth && rr.Code != http.StatusOK {
			t.Errorf("%s: esperado código de status 200, mas obteve %d", e.name, rr.Code)
		}

		// Verificar se o código de status está correto quando não autenticado
		if !e.isAuth && rr.Code != http.StatusTemporaryRedirect {
			t.Errorf("%s: esperado código de status 307, mas obteve %d", e.name, rr.Code)
		}
	}
}
