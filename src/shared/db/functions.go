package shareddb
import (
    "os"
    "fmt"
    "database/sql"
    "sort"
)
import (
    "github.com/lib/pq"
)


func buildConnStrFromEnvFn() (string, error) {
    const fn = "buildConnStrFromEnvFn"
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


func GetConn() (*sql.DB, error) {
    const wrap = "GetConn"
    // Get conn string from env
    connStr, err := buildConnStrFromEnvFn()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", wrap, err)
    }
    // Open sql conn
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to open DB: %w", wrap, err)
    }
    // Check sql conn
    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, fmt.Errorf("%s: failed to ping DB: %w", wrap, err)
    }
    return db, nil
}


func HandlePgErrorFn(table string, err error) (int, error) {
    fn := "HandlePgErrorFn"
    if pqErr, ok := err.(*pq.Error); ok {
        switch pqErr.Code {
        case "23505":// User already exist/conflict
            return 409, fmt.Errorf("%s: %s already exists: %w", fn, table, err)
        case "23514":// Invalid data/format
            return 422, fmt.Errorf("%s: invalid %s data/format: %w", fn, table, err)
        case "42703":
            return 400, fmt.Errorf("%s: unknown column used: %w", fn, err)
        default: // Failed to execute query
        return 500, fmt.Errorf("%s: failed to execute query: %w", fn, err)
        }
    }
    return 500, fmt.Errorf("%s: unexpected error: %w", fn, err)
}


func HandleSelectErrorFn(err error) (int, error) {
    fn := "HandleSelectErrorFn"
    if err == sql.ErrNoRows {
        return 404, fmt.Errorf("%s: user not found/dosen't exist", fn)
    }
    if err != nil {
        return 500, fmt.Errorf("%s: failed to execute query: %w", fn, err)
    }
    return 200, nil
}


func CheckRowsAffectedInsertFn(result sql.Result) error {
    fn := "CheckRowsAffectedInsertFn"
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: failed to check rows affected: %w", fn, err)
    }
    if rows != 1 {
        return fmt.Errorf("%s: expected 1 row affected, got %d", fn, rows)
    }
    return nil
}


func BuildSetPartsFn(data map[string]interface{}) ([]string, []interface{}, error) {
    fn := "BuildSetPartsFn"
    // Check
    if len(data) == 0 {
        return nil, nil, fmt.Errorf("%s: no fields to update", fn)
    }

    // Initialization
    setParts := []string{}
    args := []interface{}{}
    i := 1

    // Deterministic iteration
    sorted_keys := make([]string, 0, len(data))
    for k := range data {
        sorted_keys = append(sorted_keys, k)
    }
    sort.Strings(sorted_keys)
    for _, value := range sorted_keys {
        setParts = append(setParts, fmt.Sprintf("%s = $%d", value, i))
        args = append(args, data[value])
        i++
    }

    return setParts, args, nil
}

// TODO: after %s: should be lowercase letter (No rows were affected worng)
func CheckRowsAffectedUpdateFn(result sql.Result) (int, error) {
    fn := "CheckRowsAffectedUpdateFn"
    rows, err := result.RowsAffected()
    if err != nil {
        return 500, fmt.Errorf("%s: failed to check rows affected: %w", fn, err)
    }
    if rows == 0 {
        return 404, fmt.Errorf("%s: No rows were affected", fn)
    }
    return 200, nil
}



