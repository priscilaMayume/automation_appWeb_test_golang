package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// contextKey representa a chave para o contexto
type contextKey string

// contextUserKey é a chave para o contexto que armazena o IP do usuário
const contextUserKey contextKey = "user_ip"

// ipFromContext extrai o IP do contexto
func (app *application) ipFromContext(ctx context.Context) string {
	return ctx.Value(contextUserKey).(string)
}

// addIPToContext adiciona o IP ao contexto para cada solicitação
func (app *application) addIPToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		// obter o IP (o mais preciso possível)
		ip, err := getIP(r)
		if err != nil {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			if len(ip) == 0 {
				ip = "unknown"
			}
			ctx = context.WithValue(r.Context(), contextUserKey, ip)
		} else {
			ctx = context.WithValue(r.Context(), contextUserKey, ip)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getIP obtém o endereço IP do usuário a partir da solicitação HTTP
func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}

	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	if len(ip) == 0 {
		ip = "forward"
	}

	return ip, nil
}

func (app *application) auth(next http.Handler) http.Handler {
	// Retorna um manipulador HTTP que envolve a função de próxima manipulação
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica se a sessão contém a chave "user"
		if !app.Session.Exists(r.Context(), "user") {
			// Se não existir, coloca uma mensagem de erro na sessão
			app.Session.Put(r.Context(), "error", "Log in first!")
			// Redireciona o usuário para a página inicial
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// Se o usuário estiver autenticado, chama o próximo manipulador
		next.ServeHTTP(w, r)
	})
}
