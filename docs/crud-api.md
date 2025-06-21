# CRUD module

A simple API for users, events, reviews [potentially more in future].
Its purpose is to provide Create, Read, Update, and Delete access to the database  
with minimal restrictions, enforcing only data formats and models.  

It is intended to be used by a higher-level, more abstract service that will  
perform additional validation and logic as needed.

Functional Requirements:
- Provide CRUD-e API endpoints.

Non-Functional Requirements:
- Enforce username uniqueness.
- Enforce data format and type consistency (both at API and DB level).
- Allow `UPDATE` to modify one or multiple columns.

Use Case:
- Handle Create (INSERT), Read (SELECT), Update (UPDATE), and Delete (DELETE)
    operations.


## Users
<!-- {{{ Users -->
POST /create/user
payload format:
Headers:
    Content-Type: application/json
Body:
```
    {
        "username":     string  (required, mina_len: 3, max_len: 30,
                                pattern: ^[a-zA-Z0-9_]+$)
        "salt":         string  (required, hex-string, len:64)
        "hash":         string  (required, hex-string, len:64)
        "enc_symkey":   string  (required, hex-string, len:120)
    }
```

## API Responses
<!-- {{{ Responses: 201, 400, 409, 422, 500 -->
> `username` might be ommited if no username provided or wrong content-type

    201 Created
```
    {
        "message":  "Success: create user '{username}'",
        "error":    nil,
        "data":     nil,
    }
```
    400 Bad Request
```
    {
        "message":  "Fail: create user '{username}'",
        "error":    "Invalid JSON",
        "data":     nil,
    }
```
    409 Conflict 
```
    {
        "message":  "Fail: create user '{username}'",
        "error":    "User already exist",
        "data":     nil,
    }
```
    422 Unprocessable Entity
```
    {
        "message":  "Fail: create user '{username}'",
        "error":    "Invalid input format: [hash/salt/..]: [reason what is wrong]",
        "data":     nil,
    }
```
    500 Internal Server Error
```
    {
        "message":  "Fail: create user '{username}'",
        "error":    "Unknown error occured", "Internal server error"
        "data":     nil,
    }
```
<!-- Response }}} -->
## Flow
<!-- {{{ Flow -->

## Middleware

## Endpoint

## Function
### Wrapper: `InsertUser(db *sql.DB, user models.User) (int, error)`
Create query insert to database, check insertion result<br>

Requirements:
- pointer to sql.DB instance
- instance: [`User`](shared.md#struct-usermodel) from shared/models
- function: [`HandlePgError()`](shared.md#handlepgerrorerr-error-int-error) from shared/db
- function: [`CheckRowsAffectedInsert()`](shared.md#checkrowsaffectedinsertrows-int64-error) from shared/db<br>

Logic:
- Create sql query
- Call db.Exec
- Call HandlePgError()
- Call CheckRowsAffectedInsert<br>

Returns:
- `int`:    http status code
- `error`:  if execustion wasn't successful + explanation why<br>

Side effects:
Change DB, if successful insert user<br><br>
<!-- Flow }}} -->
<!-- Users }}} -->




