package main

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2" // Importa o pacote scs
)

// getSession cria e configura um gerenciador de sessão usando o pacote scs
func getSession() *scs.SessionManager {
	session := scs.New() // Cria um novo gerenciador de sessão

	// Configurações da sessão
	session.Lifetime = 24 * time.Hour // Tempo de vida da sessão (24 horas)
	session.Cookie.Persist = true // Persiste o cookie após o navegador ser fechado
	session.Cookie.SameSite = http.SameSiteLaxMode // Define o SameSite como Lax para segurança
	session.Cookie.Secure = true // Usa cookies seguros (apenas HTTPS)

	return session // Retorna o gerenciador de sessão configurado
}
