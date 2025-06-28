package cruduser
import (
    "fmt"
    "database/sql"
)
import (
    smodels "github.com/FAH2S/diar4/src/shared/models"
    sdb "github.com/FAH2S/diar4/src/shared/db"
)


func InsertUser(db *sql.DB, user smodels.User) (int, error) {
    wrap := "InsertUser"
    // Create sql query
    query := `
        INSERT INTO users (username, salt, hash, enc_symkey)
        VALUES ($1, $2, $3, $4)
    `
    // Insert
    result, err := db.Exec(query, user.Username, user.Salt, user.Hash, user.EncSymkey)
    // Map error codes to status codes
    if err != nil {
        statusCode, err := sdb.HandlePgErrorFn("user", err)
        return statusCode, fmt.Errorf("%s: %w", wrap, err)
    }
    // Check rows affected
    if err = sdb.CheckRowsAffectedInsertFn(result); err != nil {
        return 500, fmt.Errorf("%s: %w", wrap, err)
    }

    return 201, nil
}


func SelectUser(db *sql.DB, username string) (int, *smodels.User, error) {
    wrap := "SelectUser"
    // Create query
    query := `
        SELECT username, salt, hash, enc_symkey FROM users
        WHERE username = $1 LIMIT 1;
    `
    // Create user instance
    var user smodels.User
    // Query row inser + Scan load result into user
    err := db.QueryRow(query, username).Scan(
        &user.Username,
        &user.Salt,
        &user.Hash,
        &user.EncSymkey,
    )
    // Check for errors 404, 500, otherwise 200
    statusCode, err := sdb.HandleSelectErrorFn(err)
    if err != nil {
        err = fmt.Errorf("%s: %w", wrap, err)
        return statusCode, nil, err
    }
    return statusCode, &user, nil
}


