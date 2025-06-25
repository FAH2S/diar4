package integration
import (
    "testing"
    "fmt"
    "os"
    "database/sql"
    "context"
    "path/filepath"
)
import (
    "github.com/testcontainers/testcontainers-go/wait"
    "github.com/testcontainers/testcontainers-go"
)


var postgresContainer testcontainers.Container
var ctx context.Context


//{{{ helper
func getPostgresConnStr() (string, error) {
    host, err := postgresContainer.Host(ctx)
    if err != nil {
        return "", err
    }
    port, err := postgresContainer.MappedPort(ctx, "5432")
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port()), nil
}


func startPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
    initFilePath, err := filepath.Abs("../../db/init.sql")
    if err != nil {
        return nil, fmt.Errorf("Failed to get absolute path of init.sql: %w", err)
    }

    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "testuser",
            "POSTGRES_PASSWORD": "testpass",
            "POSTGRES_DB":       "testdb",
        },
        Mounts: testcontainers.Mounts(
            testcontainers.BindMount(initFilePath, "/docker-entrypoint-initdb.d/init.sql"),
        ),
        WaitingFor: wait.ForListeningPort("5432/tcp"),
    }
    return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
}


func checkPostgresConnection() error {
    connStr, err := getPostgresConnStr()

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to open DB: %w", err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        return fmt.Errorf("failed to ping DB: %w", err)
    }

    return nil

}
//}}} helper


func TestMain(m *testing.M) {
    ctx = context.Background()

    var err error
    postgresContainer, err = startPostgresContainer(ctx)
    if err != nil {
        panic(fmt.Errorf("Failed to start postgres container: %w", err))
    }
    if err := checkPostgresConnection(); err != nil {
        panic(fmt.Errorf("DB connection test failed: %w", err))
    }

    code := m.Run()

    _ = postgresContainer.Terminate(ctx)

    os.Exit(code)
}


