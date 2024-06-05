package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// O utilizador descreve os dados para o tipo de utilizador.
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsAdmin   int       `json:"is_admin"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// PasswordMatches usa o pacote bcrypt do Go para comparar uma senha fornecida pelo usuário
// com o hash que temos armazenado para um determinado usuário no banco de dados. Se a senha
// e o hash coincidirem, retornamos true; caso contrário, retornamos false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// senha invalida
			return false, nil
		default:  
			return false, err
		}
	}

	return true, nil
}
