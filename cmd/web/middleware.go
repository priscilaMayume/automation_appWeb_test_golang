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
				ip = "desconhecido"
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
		return "desconhecido", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q não é IP:porta", r.RemoteAddr)
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
