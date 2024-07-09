package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"
)

const jwtTokenExpiry = time.Minute * 15
const refreshTokenExpiry = time.Hour * 24

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

// Função para obter o token do cabeçalho da requisição e verificar sua validade
func (app *application) getTokenFromHeaderandVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// Esperamos que o cabeçalho de autorização tenha o seguinte formato:
	// Bearer <token>
	// adicionar um cabeçalho 
	w.Header().Add("Vary", "Authorization")

	// obter o cabeçalho de autorização
	authHeader := r.Header.Get("Authorization")

	// verificação básica
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	// dividir o cabeçalho nos espaços
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	// verificar se temos a palavra "Bearer"
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("unauthorized: no Bearer")
	}

	token := headerParts[1]

	// declarar uma variável Claims vazia
	claims := &Claims{}

	// analisar o token com nossas claims (lendo para claims), usando nosso segredo (do receiver)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// validar o algoritmo de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.JWTSecret), nil
	})

	// verificar se há um erro; note que isso também captura tokens expirados.
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	// garantir que fomos nós que emitimos este token
	if claims.Issuer != app.Domain {
		return "", nil, errors.New("incorrect issuer")
	}

	// token válido
	return token, claims, nil
}

func (app *application) generateTokenPair(user *data.User) (TokenPairs, error) {
	// Create the token.
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName) // Define o nome completo do usuário nos claims.
	claims["sub"] = fmt.Sprint(user.ID) // Define o ID do usuário nos claims.
	claims["aud"] = app.Domain // Define o domínio para o qual o token é emitido nos claims.
	claims["iss"] = app.Domain // Define o emissor (issuer) do token nos claims.
	if user.IsAdmin == 1 {
		claims["admin"] = true // Define se o usuário é um administrador nos claims.
	} else {
		claims["admin"] = false // Define que o usuário não é um administrador nos claims.
	}

	// set the expiry
	claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix() // Define o tempo de expiração do token nos claims.

	// create the signed token
	signedAccessToken, err := token.SignedString([]byte(app.JWTSecret)) // Gera o token de acesso assinado.
	if err != nil {
		return TokenPairs{}, err
	}

	// create the refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID) // Define novamente o ID do usuário nos claims do token de atualização.
	
	// set expiry; must be longer than jwt expiry
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix() // Define o tempo de expiração do token de atualização.

	// create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(app.JWTSecret)) // Gera o token de atualização assinado.
	if err != nil {
		return TokenPairs{}, err
	}

	var tokenPairs = TokenPairs{
		Token: signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokenPairs, nil
}
