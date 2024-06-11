package dbrepo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"
)

type TestDBRepo struct{}

func (m *TestDBRepo) Connection() *sql.DB {
	return nil
}

// AllUsers retorna todos os usuários como um slice de *data.User
func (m *TestDBRepo) AllUsers() ([]*data.User, error) {
	var users []*data.User

	return users, nil
}

// GetUser retorna um usuário pelo id
func (m *TestDBRepo) GetUser(id int) (*data.User, error) {
	var user = data.User{
		ID: 1,
	}

	return &user, nil
}

// GetUserByEmail retorna um usuário pelo endereço de e-mail
func (m *TestDBRepo) GetUserByEmail(email string) (*data.User, error) {
	if email == "admin@example.com" {
		user := data.User{
			ID: 1,
			FirstName: "Admin",
			LastName: "User",
			Email: "admin@example.com",
			Password: "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			IsAdmin: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return &user, nil
	}
	return nil, errors.New("not found")
}

// UpdateUser atualiza um usuário no banco de dados
func (m *TestDBRepo) UpdateUser(u data.User) error {
	return nil
}

// DeleteUser exclui um usuário do banco de dados pelo id
func (m *TestDBRepo) DeleteUser(id int) error {
	return nil
}

// InsertUser insere um novo usuário no banco de dados e retorna o ID da linha recém-inserida
func (m *TestDBRepo) InsertUser(user data.User) (int, error) {
	return 2, nil
}

// ResetPassword é o método que usaremos para mudar a senha de um usuário.
func (m *TestDBRepo) ResetPassword(id int, password string) error {
	return nil
}

// InsertUserImage insere uma imagem de perfil de usuário no banco de dados.
func (m *TestDBRepo) InsertUserImage(i data.UserImage) (int, error) {
	return 1, nil
}
