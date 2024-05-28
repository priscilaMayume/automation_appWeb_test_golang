package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_addIPToContext(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false},
		{"", "", "hello:world", false},
	}

	var app application

	// cria um handler dummy que usaremos para verificar o contexto
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verifica se o valor existe no contexto
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Error(contextUserKey, "não presente")
		}

		// certifica-se de que recebemos uma string de volta
		ip, ok := val.(string)
		if !ok {
			t.Error("não é string")
		}
		t.Log(ip)
	})

	for _, e := range tests {
		// cria o handler para testar
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
	// criar um app var de tipo application
	var app application

	// obter um contexto
	ctx := context.Background()

	// colocar algo no contexto
	ctx = context.WithValue(ctx, contextUserKey, "whatever")

	// chamar a função
	ip := app.ipFromContext(ctx)

	// peformance do teste
	if !strings.EqualFold("whatever", ip) {
		t.Error("wrong value returned from context")
	}
}