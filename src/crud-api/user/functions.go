package user
import (
    "fmt"
    "database/sql"
)
import (
    smodels "github.com/FAH2S/diar4/src/shared/models"
    sdb "github.com/FAH2S/diar4/src/shared/db"
)


func InsertUser(db *sql.DB, user smodels.User) (int, error) {
    fn := "insertUser"
    // Create sql query
    query := `
        INSERT INTO users (username, salt, hash, enc_symkey)
        VALUES ($1, $2, $3, $4)
    `
    // Insert
    result, err := db.Exec(query, user.Username, user.Salt, user.Hash, user.EncSymkey)
    // Map error codes to status codes
    if err != nil {
        statusCode, err := sdb.HandlePgError(err)
        return statusCode, fmt.Errorf("%s: %w", fn, err)
    }
    // Check rows affected
    if err = sdb.CheckRowsAffectedInsert(result); err != nil {
        return 500, fmt.Errorf("%s: %w", fn, err)
    }

    return 201, nil
}


