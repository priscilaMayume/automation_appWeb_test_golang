package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
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
