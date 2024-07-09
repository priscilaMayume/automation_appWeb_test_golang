package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Função para abrir a conexão com o banco de dados usando o DSN fornecido
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Verificando a conexão com o banco de dados
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Método para a aplicação conectar ao banco de dados
func (app *application) connectToDB() (*sql.DB, error) {
	connection, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres!")

	return connection, nil
}
