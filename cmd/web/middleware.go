package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// Define um tipo contextKey como string
type contextKey string

// Define uma constante contextUserKey para armazenar o IP do usuário no contexto
const contextUserKey contextKey = "user_ip"

// ipFromContext recupera o IP do usuário a partir do contexto.
func (app *application) ipFromContext(ctx context.Context) string {
	return ctx.Value(contextUserKey).(string)
}

// addIPToContext é um middleware que adiciona o IP do usuário ao contexto da solicitação.
func (app *application) addIPToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()

		// Obtém o IP (da forma mais precisa possível)
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

		// Passa o contexto atualizado para o próximo handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getIP obtém o IP do cliente a partir do cabeçalho da solicitação ou do endereço remoto.
func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown", err
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("userip: %q is not IP:port", r.RemoteAddr)
	}

	// Verifica o cabeçalho "X-Forwarded-For" para IPs encaminhados
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	if len(ip) == 0 {
		ip = "forward"
	}

	return ip, nil
}
