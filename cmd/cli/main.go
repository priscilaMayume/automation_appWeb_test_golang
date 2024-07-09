package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Estrutura da aplicação contendo as configurações necessárias.
type application struct {
	JWTSecret string // Segredo JWT para assinatura do token.
	Action    string // Ação a ser realizada: 'valid' para token válido, 'expired' para token expirado.
}

// Função principal responsável por gerar um token JWT com base nas configurações definidas.
func main() {
	var app application

	// Definindo as flags de linha de comando para segredo JWT e ação (valid/expired).
	flag.StringVar(&app.JWTSecret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "segredo")
	flag.StringVar(&app.Action, "action", "valid", "ação: valid|expired")
	flag.Parse()

	// Gerando um novo token JWT.
	token := jwt.New(jwt.SigningMethodHS256)

	// Definindo as claims do token.
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "John Doe"      // Nome do usuário associado ao token.
	claims["sub"] = "1"              // Identificador único do usuário.
	claims["admin"] = true           // Indicação se o usuário é administrador.
	claims["aud"] = "example.com"    // Audiência do token.
	claims["iss"] = "example.com"    // Emissor do token.

	// Definindo a expiração do token.
	if app.Action == "valid" {
		expires := time.Now().UTC().Add(time.Hour * 72) // Token válido por 72 horas.
		claims["exp"] = expires.Unix()
	} else {
		expires := time.Now().UTC().Add(time.Hour * 100 * -1) // Token expirado (hora atual menos 100 horas).
		claims["exp"] = expires.Unix()
	}

	// Criando o token como uma string assinada.
	if app.Action == "valid" {
		fmt.Println("VALID Token:")
	} else {
		fmt.Println("EXPIRED Token:")
	}
	signedAccessToken, err := token.SignedString([]byte(app.JWTSecret))
	if err != nil {
		log.Fatal(err)
	}
	// Imprimindo o token gerado no console.
	fmt.Println(string(signedAccessToken))
}
