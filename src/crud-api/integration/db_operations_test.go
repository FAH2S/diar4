package integration
import (
    "context"
    "testing"
    "fmt"
    "os"
    "database/sql"
    "path/filepath"
    "errors"
    "strings"
)
import (
    _ "github.com/lib/pq"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    smodels "github.com/FAH2S/diar4/src/shared/models"
    cruduser "github.com/FAH2S/diar4/src/crud-api/user"
)


var postgresContainer testcontainers.Container
var ctx context.Context


//{{{ DRY
func getPostgresConnStr(t *testing.T) string {
    host, err := postgresContainer.Host(ctx)
    if err != nil {
        t.Fatal(err)
    }
    port, err := postgresContainer.MappedPort(ctx, "5432")
    if err != nil {
        t.Fatal(err)
    }

    return fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
}
//}}} DRY


//{{{ Setup postgres docker
//TODO: potentially extract as fn for later testing
func TestMain(m *testing.M) {
    ctx = context.Background()
    initFilePath, err := filepath.Abs("../../db/init.sql")
    if err != nil {
        panic(fmt.Errorf("Failed to get absolute path of init.sql: %w", err))
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

    postgresContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        panic(fmt.Errorf("failed to start postgres container: %w", err))
    }

    code := m.Run()

    _ = postgresContainer.Terminate(ctx)

    os.Exit(code)
}


func TestPostgresConnection(t *testing.T) {
    connStr := getPostgresConnStr(t)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        t.Fatal("failed to open DB:", err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        t.Fatal("failed to ping DB:", err)
    }

    t.Log("Successfully connected to postgres container")
}
//}}} Setup postgres docker


//{{{ Insert user
func TestInsertUser(t *testing.T) {
    connStr := getPostgresConnStr(t)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        t.Fatal("failed to open DB:", err)
    }
    defer db.Close()

    tests := []struct {
        name                string
        user                smodels.User
        expectedStatusCode  int
        expectedError       error
    }{
        {
            name:               "validInput",
            user:               smodels.User{
                Username:   "test_user_123",
                Salt:       "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 201,
            expectedError:      nil,
        }, {
            name:               "userAlreadyExists",
            user:               smodels.User{
                Username:   "test_user_123",
                Salt:       "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 409,
            expectedError:      errors.New("user already exists"),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableSalt",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       "344feecf40d375380ed5fe0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_salt_check\""),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableHash",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       "344feecf40d375380ed5fe0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a7$b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_hash_check\""),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableEnc_symkey",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       "344feecf40d375380ed5fe0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a7$b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_enc_symkey_check\""),
        },
    }

    for _, tc := range tests{
        t.Run(tc.name, func(t *testing.T) {
            actual, err := cruduser.InsertUser(db, tc.user)
            if tc.expectedStatusCode != actual {
                t.Errorf("\nExpected:\t%d\nGot:\t%d", tc.expectedStatusCode, actual)
            }

            if (err == nil) != (tc.expectedError == nil) {
                t.Errorf("\nExpected:\t%v\nGot:\t%v", tc.expectedError, err)
            } else if err != nil &&
                tc.expectedError != nil &&
                !strings.Contains(err.Error(), tc.expectedError.Error()){
                t.Errorf("\nExpected to contain:\t%q\nGot:\t\t\t%q", tc.expectedError, err)
            }
        })
    }
}
//}}} Insert user


