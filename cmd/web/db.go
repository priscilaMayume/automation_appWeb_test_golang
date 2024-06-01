package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// openDB abre uma conexão com o banco de dados usando a string de conexão fornecida (DSN).
func openDB(dsn string) (*sql.DB, error) {
	// Abre a conexão com o banco de dados usando o driver "pgx".
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Verifica se a conexão com o banco de dados está ativa.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Retorna a conexão do banco de dados.
	return db, nil
}

// connectToDB conecta ao banco de dados usando a DSN armazenada na estrutura da aplicação.
func (app *application) connectToDB() (*sql.DB, error) {
	// Abre a conexão com o banco de dados.
	connection, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}

	// Loga uma mensagem indicando que a conexão foi bem-sucedida.
	log.Println("Connected to Postgres!")

	// Retorna a conexão do banco de dados.
	return connection, nil
}