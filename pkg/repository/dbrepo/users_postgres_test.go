//go:build integration

package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/data"
	"github.com/priscilaMayume/automation_appWeb_test_golang/pkg/repository"
)

var (
	host     = "localhost"  // Host do banco de dados
	user     = "postgres"   // Usuário do banco de dados
	password = "postgres"   // Senha do banco de dados
	dbName   = "users_test" // Nome do banco de dados de teste
	port     = "5435"       // Porta do banco de dados
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5" // String de conexão com o banco de dados
	
)

var resource *dockertest.Resource // Recurso do Docker para o banco de dados
var pool *dockertest.Pool         // Pool de conexões do Docker
var testDB *sql.DB                // Instância do banco de dados de teste
var testRepo repository.DatabaseRepo

// TestMain é a função principal de testes
// Ela é executada antes de qualquer outro teste
func TestMain(m *testing.M) {
	// Conecta ao Docker; falha se o Docker não estiver em execução
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	// Configura as opções do Docker, especificando a imagem e assim por diante
	opts := dockertest.RunOptions{
		Repository: "postgres", // Repositório da imagem Docker
		Tag: "14.5", // Tag da imagem Docker
		Env: []string{ // Variáveis de ambiente para o container Docker
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"}, // Porta exposta pelo container
		PortBindings: map[docker.Port][]docker.PortBinding{ // Mapeamento de portas
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// Obtém um recurso (imagem Docker)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// Inicia a imagem e espera até que esteja pronta
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// Popula o banco de dados com tabelas vazias (esta parte pode ser personalizada conforme necessário)
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = &PostgresDBRepo{DB: testDB}

	// Executa os testes
	code := m.Run()

	// Limpa os recursos (encerra e remove o container Docker)
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	// Sai do programa com o código de status retornado pelos testes
	os.Exit(code)
}

// createTables lê o arquivo SQL contendo as instruções para criar tabelas
// e executa essas instruções no banco de dados.
func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Test_pingDB verifica se é possível fazer ping no banco de dados.
// Se o ping falhar, o teste registra um erro.
func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

// Testa a inserção de um usuário no repositório Postgres
func TestPostgresDBRepoInsertUser(t *testing.T) {
	// Cria um usuário de teste com dados fictícios
	testUser := data.User{
		FirstName: "Admin",
		LastName: "User",
		Email: "admin@example.com",
		Password: "secret",
		IsAdmin: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insere o usuário de teste no repositório e verifica se ocorreu algum erro
	id, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user returned an error: %s", err)
	}

	// Verifica se o ID retornado é 1, caso contrário, reporta um erro
	if id != 1 {
		t.Errorf("insert user returned wrong id; expected 1, but got %d", id)
	}
}

// Testa a obtenção de todos os usuários do repositório Postgres
func TestPostgresDBRepoAllUsers(t *testing.T) {
	// Obtém todos os usuários do repositório e verifica se houve algum erro
	users, err := testRepo.AllUsers()
	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	// Verifica se o número de usuários retornado é 1, caso contrário, reporta um erro
	if len(users) != 1 {
		t.Errorf("all users reports wrong size; expected 1, but got %d", len(users))
	}

	// Cria um segundo usuário de teste com dados fictícios
	testUser := data.User{
		FirstName: "Jack",
		LastName: "Smith",
		Email: "jack@smith.com",
		Password: "secret",
		IsAdmin: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insere o segundo usuário de teste no repositório
	_, _ = testRepo.InsertUser(testUser)

	// Obtém todos os usuários novamente e verifica se houve algum erro
	users, err = testRepo.AllUsers()
	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	// Verifica se o número de usuários retornado é 2 após a inserção, caso contrário, reporta um erro
	if len(users) != 2 {
		t.Errorf("all users reports wrong size after insert; expected 2, but got %d", len(users))
	}
}

// Testa a obtenção de um usuário pelo ID no repositório Postgres
func TestPostgresDBRepoGetUser(t *testing.T){
	// Obtém um usuário pelo ID 1 e verifica se houve algum erro
	user, err := testRepo.GetUser(1)
	if err != nil {
		t.Errorf("error getting user by id: %s", err)
	}

	// Verifica se o email do usuário retornado é "admin@example.com", caso contrário, reporta um erro
	if user.Email != "admin@example.com" {
		t.Errorf("wrong email returned by GetUser; expected admin@example.com but got %s", user.Email)
	}

	// Tenta obter um usuário inexistente pelo ID 3 e verifica se houve algum erro
	_, err = testRepo.GetUser(3)
	if err == nil {
		t.Error("no error reported when getting non existent user by id")
	}
}

// Testa a obtenção de um usuário pelo email no repositório Postgres
func TestPostgresDBRepoGetUserByEmail(t *testing.T){
	// Obtém um usuário pelo email "jack@smith.com" e verifica se houve algum erro
	user, err := testRepo.GetUserByEmail("jack@smith.com")
	if err != nil {
		t.Errorf("error getting user by email: %s", err)
	}

	// Verifica se o ID do usuário retornado é 2, caso contrário, reporta um erro
	if user.ID != 2 {
		t.Errorf("wrong id returned by GetUserByEmail; expected 2 but got %d", user.ID)
	}
}

// Testa a redefinição da senha de um usuário no repositório Postgres
func TestPostgresDBRepoResetPassword(t *testing.T) {
	// Redefine a senha do usuário com ID 1 para "password" e verifica se houve algum erro
	err := testRepo.ResetPassword(1, "password")
	if err != nil {
		t.Error("error resetting user's password", err)
	}

	// Obtém o usuário com ID 1 para verificar se a senha foi redefinida corretamente
	user, _ := testRepo.GetUser(1)
	// Verifica se a nova senha corresponde à senha redefinida
	matches, err := user.PasswordMatches("password")
	if err != nil {
		t.Error(err)
	}

	// Se a senha não corresponder, reporta um erro
	if !matches {
		t.Errorf("password should match 'password', but does not")
	}
}

// Testa a inserção de uma imagem de usuário no repositório Postgres
func TestPostgresDBRepoInsertUserImage(t *testing.T) {
	// Cria uma imagem de usuário de teste com dados fictícios
	var image data.UserImage
	image.UserID = 1
	image.FileName = "test.jpg"
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	// Insere a imagem do usuário no repositório e verifica se houve algum erro
	newID, err := testRepo.InsertUserImage(image)
	if err != nil {
		t.Error("inserting user image failed:", err)
	}

	// Verifica se o ID retornado é 1, caso contrário, reporta um erro
	if newID != 1 {
		t.Error("got wrong id for image; should be 1, but got", newID)
	}

	// Tenta inserir uma imagem de usuário com um UserID inexistente e verifica se houve algum erro
	image.UserID = 100
	_, err = testRepo.InsertUserImage(image)
	if err == nil {
		t.Error("inserted a user image with non-existent user id")
	}
}
