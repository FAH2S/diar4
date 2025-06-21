package shareddb
import (
    "os"
    "fmt"
    "database/sql"
)
import (
    "github.com/lib/pq"
)

func buildConnStrFromEnv() (string, error) {
    const fn = "buildConnStrFromEnv"
    keys := []string{"DB_USER", "DB_PWD", "DB_HOST", "DB_PORT", "DB_NAME"}
    values := make(map[string]string)

    // Iterate over list, fetch values update dict/map, raise error if empty
    for _, key := range keys{
        val := os.Getenv(key)
        if val == "" {
            return "", fmt.Errorf("%s: Required environment variable %s, not set\n", fn, key)
        }
        values[key] = val
    }

    // Construct conn string
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        values["DB_USER"],
        values["DB_PWD"],
        values["DB_HOST"],
        values["DB_PORT"],
        values["DB_NAME"],
    )
    //fmt.Println(connStr)
    return connStr, nil
}


func HandlePgError(err error) (int, error) {
    fn := "HandlePgError"
    if pqErr, ok := err.(*pq.Error); ok {
        switch pqErr.Code {
        case "23505":// User already exist/conflict
            return 409, fmt.Errorf("%s: user already exists: %w", fn, err)
        case "23514":// Invalid data/format
            return 422, fmt.Errorf("%s: invalid user data/format: %w", fn, err)
        case "42703":
            return 400, fmt.Errorf("%s: unknown column used: %w", fn, err)
        default: // Failed to execute query
        return 500, fmt.Errorf("%s: failed to execute query: %w", fn, err)
        }
    }
    return 500, fmt.Errorf("%s: unexpected error: %w", fn, err)
}


func CheckRowsAffectedInsert(result sql.Result) error {
    fn := "CheckRowsAffectedInsert"
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: failed to check rows affected: %w", fn, err)
    }
    if rows != 1 {
        return fmt.Errorf("%s: expected 1 row affected, got %d", fn, rows)
    }
    return nil
}


