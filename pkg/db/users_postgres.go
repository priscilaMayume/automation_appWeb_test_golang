package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

type PostgresConn struct {
	DB *sql.DB
}

// AllUsers devolve todos os utilizadores como uma fatia de *data.User
func (m *PostgresConn) AllUsers() ([]*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, is_admin, created_at, updated_at
	from users order by last_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*data.User

	for rows.Next() {
		var user data.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetUser devolve um utilizador por id
func (m *PostgresConn) GetUser(id int) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, email, first_name, last_name, password, is_admin, created_at, updated_at 
		from 
			users 
		where 
		    id = $1`

	var user data.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail devolve um utilizador por endereço de correio eletrónico
func (m *PostgresConn) GetUserByEmail(email string) (*data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			id, email, first_name, last_name, password, is_admin, created_at, updated_at 
		from 
			users 
		where 
		    email = $1`

	var user data.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser actualiza um utilizador na base de dados
func (m *PostgresConn) UpdateUser(u data.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		is_admin = $4,
		updated_at = $5
		where id = $6
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.IsAdmin,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser apaga um utilizador da base de dados, por id
func (m *PostgresConn) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// InsertUser insere um novo utilizador na base de dados e devolve o ID da nova linha inserida
func (m *PostgresConn) InsertUser(user data.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, is_admin, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = m.DB.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.IsAdmin,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword é o método que vamos utilizar para alterar a palavra-passe de um utilizador.
func (m *PostgresConn) ResetPassword(id int, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = m.DB.ExecContext(ctx, stmt, hashedPassword, id)
	if err != nil {
		return err
	}

	return nil
}

// InsertUserImage insere uma imagem de perfil de utilizador na base de dados.
func (m *PostgresConn) InsertUserImage(i data.UserImage) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into user_images (user_id, fileName, created_at, updated_at)
		values ($1, $2, $3, $4) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		i.UserID,
		i.FileName,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}
