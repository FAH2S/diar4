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
    sdb "github.com/FAH2S/diar4/src/shared/db"
)


var postgresContainer testcontainers.Container
var ctx context.Context
var db *sql.DB


//{{{ helper
func initialzeDbTestEnv() error {
    host, err := postgresContainer.Host(ctx)
    if err != nil {
        return err
    }
    port, err := postgresContainer.MappedPort(ctx, "5432")
    if err != nil {
        return err
    }
    os.Clearenv()
    os.Setenv("DB_USER", "testuser")
    os.Setenv("DB_PWD", "testpass")
    os.Setenv("DB_NAME", "testdb")
    os.Setenv("DB_HOST", host)
    os.Setenv("DB_PORT", port.Port())

    return nil
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
    // testcontainers.GenericContainer retrurns container + error on its own no need 
    //  for explicit error(nil) return
    return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
}
//}}} helper


func TestMain(m *testing.M) {
    ctx = context.Background()

    //need to declare err becuse we can't use := operator on globar var
    //  in this case `postgresContainer`
    var err error
    postgresContainer, err = startPostgresContainer(ctx)
    if err != nil {
        panic(fmt.Errorf("Failed to start postgres container: %w", err))
    }
    err = initialzeDbTestEnv()
    if err != nil {
        panic(fmt.Errorf("Failed to initialize DB env's"))
    }
    db, err = sdb.GetConn()
    if err != nil {
        panic(fmt.Errorf("Failed to initialize db conn: %w", err))
    }
    // Run tests
    code := m.Run()
    // Teardown
    db.Close()
    _ = postgresContainer.Terminate(ctx)
    os.Exit(code)
}


