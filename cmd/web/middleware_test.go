package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
