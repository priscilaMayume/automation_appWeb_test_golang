package main

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Lê o payload JSON.
	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Busca o usuário pelo endereço de e-mail.
	user, err := app.DB.GetUserByEmail(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Verifica a senha.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Gera tokens.
	tokenPairs, err := app.generateTokenPair(user)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Envia o token para o usuário.
	_ = app.writeJSON(w, http.StatusOK, tokenPairs)
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	// Método ainda não implementado.
}

func (app *application) allUsers(w http.ResponseWriter, r *http.Request) {
	// Método ainda não implementado.
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	// Método ainda não implementado.
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	// Método ainda não implementado.
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Método ainda não implementado.
}

func (app *application) insertUser(w http.ResponseWriter, r *http.Request) {
	// Método ainda não implementado.
}
