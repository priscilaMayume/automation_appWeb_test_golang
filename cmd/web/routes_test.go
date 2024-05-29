package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// Test_application_routes testa se todas as rotas estão registradas corretamente
func Test_application_routes(t *testing.T) {
	var registered = []struct{
		route string
		method string
	}{
		{"/", "GET"},
		{"/login", "POST"},
		{"/static/*", "GET"},
	}

	mux := app.routes() // Obtém o roteador

	chiRoutes := mux.(chi.Routes) // Converte o roteador para o tipo chi.Routes

	for _, route := range registered {
		// Verifica se a rota existe
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("rota %s não está registrada", route.route)
		}
	}
}

// routeExists verifica se uma rota está registrada no roteador
func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})

	return found
}
