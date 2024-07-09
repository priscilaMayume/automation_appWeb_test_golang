package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"

	"github.com/golang-jwt/jwt/v4"
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

func (app *application) getTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// esperamos que nosso cabeçalho de autorização se pareça com isto:
	// Bearer <token>
	// adicionar um cabeçalho
	w.Header().Add("Vary", "Authorization")

	// obter o cabeçalho de autorização
	authHeader := r.Header.Get("Authorization")

	// verificação de sanidade
	if authHeader == "" {
		return "", nil, errors.New("sem cabeçalho de autorização")
	}

	// dividir o cabeçalho em espaços
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("cabeçalho de autorização inválido")
	}

	// verificar se temos a palavra "Bearer"
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("não autorizado: sem Bearer")
	}

	token := headerParts[1]

	// declarar uma variável Claims vazia
	claims := &Claims{}

	// analisar o token com nossas claims (lendo em claims), usando nosso segredo (do receiver)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// validar o algoritmo de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(app.JWTSecret), nil
	})

	// verificar se há um erro; note que isso também captura tokens expirados.
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("token expirado")
		}
		return "", nil, err
	}

	// garantir que nós emitimos este token
	if claims.Issuer != app.Domain {
		return "", nil, errors.New("emissor incorreto")
	}

	// token válido
	return token, claims, nil
}

func (app *application) generateTokenPair(user *data.User) (TokenPairs, error) {
	// Criar o token.
	token := jwt.New(jwt.SigningMethodHS256)

	// definir claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = app.Domain
	claims["iss"] = app.Domain
	if user.IsAdmin == 1 {
		claims["admin"] = true
	} else {
		claims["admin"] = false
	}

	// definir a expiração
	claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix()

	// criar o token assinado
	signedAccessToken, err := token.SignedString([]byte(app.JWTSecret))
	if err != nil {
		return TokenPairs{}, err
	}

	// criar o token de atualização
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)

	// definir a expiração; deve ser maior que a expiração do jwt
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	// criar o token de atualização assinado
	signedRefreshToken, err := refreshToken.SignedString([]byte(app.JWTSecret))
	if err != nil {
		return TokenPairs{}, err
	}

	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	return tokenPairs, nil
}
