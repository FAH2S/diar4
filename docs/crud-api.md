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


## Middleware
<!-- {{{ Middleware -->
Middleware that checks for method (POST) and header (Content-Type: application/json).
## API Respnse
    400 Bad Request
```
    {
        "message":  "Fail: to process 'URL path:[/create/user, /read/user ...]'",
        "error":    "Method not allowed"/"Content-Type must be application/json",
        "data":     nil,
    }
```
## Endpoint
### Wrapper: `ValidateMethodAndTypeEndpoint(next http.Handler) http.Handler`
Intercepts packet and check if it uses correct method and header.<br>

Requirements:
- http.Handler
- function: [`isMethodPOSTFn`](crud-api.md#function-ismethodpostfnr-httprequest-bool)
- function: [`isHeaderCTAJFn`](crud-api.md#function-isheaderctajfnr-httprequest-bool)<br>

Logic:
- Call `isMethodPOSTFn`
- Call `isHeaderCTAJFn`
- Pass packet to CRUD<br>

Retruns:
- `http.Handler`: if both check pass request is forwarded to next handler<br>
(e.g., handler for /create/user)

### Function: `isMethodPOSTFn(r *http.Request) bool`
Checks if method is POST
### Function: `isHeaderCTAJFn(r *http.Request) bool`
Check if it contains header `Content-Type: application/json`
<!-- }}} Middleware --><br>

## Users
<!-- {{{ Users -->
<!-- {{{ CREATE User -->
POST /create/user<br>
Headers:
```
    Content-Type: application/json
```
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

<!-- {{{ Responses: 201, 400, 409, 422, 500 -->
## API Responses
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
<!-- {{{ Flow -->
## Flow
## Endpoint
### Wrapper: `CreateUserEndpoint(w http.ResponseWriter, r *http.Request, db *sql.DB)`
Accept package, convert to user model, insert to DB.<br>

Requirements:
- pointer to sql.DB instance
- function: [`Validate`](shared.md#wrapper-validate-error) from shared/models
- wrapper:  [`InsertUser`](#wrapper-insertuserdb-sqldb-user-modelsuser-int-error)
- function: [`WriteJSONResponseFn`](shared.md#function-writejsonresponsefnw-httpresponsewrite-statuscode-int-message-string-errmsg-string-data-interface) from shared/api<br>

Logic:
- Decode JSON into `user` model
- Call `user.Validate()`
- Call `InsertUser()`
- return `APIResponse`<br>

Returns:
- api response [`APIResponse`](shared.md#struct-apiresponse)<br>

Side effects:
Change DB, if successful insert user (indirectly via InsertUser)<br><br>

## Function
### Wrapper: `InsertUser(db *sql.DB, user models.User) (int, error)`
Create query insert to database, check insertion result<br>

Requirements:
- pointer to sql.DB instance
- instance: [`User`](shared.md#struct-user) from shared/models
- function: [`HandlePgErrorFn()`](shared.md#function-handlepgerrorfnerr-error-int-error) from shared/db
- function: [`CheckRowsAffectedInsertFn()`](shared.md#function-checkrowsaffectedinsertfnresult-sqlresult-error) from shared/db<br>

Logic:
- Create sql query
- Call `db.Exec()`
- Call `HandlePgErrorFn()` if error present
- Call `CheckRowsAffectedInsertFn()`<br>

Returns:
- `int`:    http status code
- `error`:  if execution wasn't successful + explanation why<br>

Side effects:
Change DB, if successful insert user<br><br>
<!-- Flow }}} -->
<!-- }}}CREATE User -->


<!-- {{{ READ User -->
POST /read/user<br>
Headers:
```
    Content-Type: application/json
```
Body:
```
    {
        "username":     string  (required, mina_len: 3, max_len: 30,
                                pattern: ^[a-zA-Z0-9_]+$)
    }
```
<!-- {{{ Responses: 200, 400, 404, 422, 500 -->
## API Responses
    200 OK
```
    {
        "message":  "Success: read user '{username}'",
        "error":    nil,
        "data":     {JSON map of User model},
    }
```
    400 Bad Request
```
    {
        "message":  "Fail: read user ''",   //Malformed JSON can't process body aka extract username
        "error":    "Invalid JSON",
        "data":     nil,
    }
```
    404 Not Found
```
    {
        "message":  "Fail: read user '{username}'",
        "error":    "User not found, dosen't exist",
        "data":     nil,
    }
```
    422 Unprocessable Entity
```
    {
        "message":  "Fail: read user '{username}'",
        "error":    "Invalid input format: username: [reason what is wrong]",
        "data":     nil,
    }
```
    500 Internal Server Error
```
    {
        "message":  "Fail: read user '{username}'",
        "error":    "Unknown error occured", "Internal server error"
        "data":     nil,
    }
```
<!-- }}} Responses: 200, 400, 404, 422, 500 -->
<!-- {{{ Flow -->
## Flow
## Endpoint
### Wrapper: `ReadUserEndpoint(w http.ResponseWriter, r *http.Request, db *sql.DB)`
Accept package, checks if username is valid, fetches from DB.<br>

Requirements:
- pointer to `sql.DB` instance
- function: [`ExtractJSONValueFn()`](shared.md#function-extractjsonvaluefnr-httprequest-key-string-target-interface-error) from shared/api
- function: [`IsValidUsernameFn()`](shared.md#function-isvalidusernamefnusername-string-error) from shared/models
- wrapper:  [`SelectUser()`](#wrapper-selectuserdb-sqldb-username-string-int-modelsuser-error)
- function: [`WriteJSONResponseFn()`](shared.md#function-writejsonresponsefnw-httpresponsewrite-statuscode-int-message-string-errmsg-string-data-interface) from shared/api<br>

Logic:
- Call `ExtractJSONValueFn()` to get username
- Call `IsValidUsernameFn()`
- Call `SelectUser()`
- return `APIResponse`<br><br>

Returns:
- api response [`APIResponse`](shared.md#struct-apiresponse)<br><br>


### Wrapper: `SelectUser(db *sql.DB, username string) (int, *models.User, error)`
Create query to select/fetch user from database, check, selection result<br>

Requirements:
- pointer to `sql.DB` instance
- instance: [`User`](shared.md#struct-user) from shared/models
- wrapper: [`HandleSelectErrorFn()`](shared.md#wrapper-handleselecterrorfnerr-error-int-error) from shared/db<br>

Logic:
- Create sql query
- Create `user` instance
- Call `db.QueryRow()`, then via `.Scan()` load result into `user` instance
- Call `HandleSelectError()` <br>

Returns:
- `int`:            http status code
- `models.User`:    pointer of selected/fetched user
- `erorr`:          if execution wasn't successful + explanation why<br><br>
<!-- }}} Flow -->
<!-- }}} READ User -->


<!-- {{{ UPDATE User -->
<!-- }}} UPDATE User -->


<!-- {{{ DELETE User -->
<!-- }}} DELETE User -->
<!-- Users }}} -->




