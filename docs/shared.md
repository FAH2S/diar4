## DB
<!-- {{{ DB -->
### Wrapper: `db.GetConn() (*sql.DB, error)`
Creates new connection to database from environment variables.<br>

Requirements:
- `DB_USER, DB_PWD, DB_HOST, DB_PORT, DB_NAME` must be set as non empty string.
- function: [`buildConnStrFromEnvFn()`](shared.md#function-buildconnstrfromenvfn)
- `sql.Open()`
- `sql.Ping()`<br>

Logic:
- Call `buildCononStrFromEnvFn()` to create connection string
- Call `sql.Open()` to create DB connection
- Check sql connection via `sql.Ping()`
- Return connection or error<br>

Returns:
- `*sql.db`:    pointer to db connection
- `error`:      if any variable is missing or `sql.open` fails<br><br>


### Function: `buildConnStrFromEnvFn() (string, error)`
Builds connection string from environment variables.<br>

Logic:
- Read enviroment variables + validate
- Format connection string from enviroment variables<br>

Returns:
- `string`: connection string ex.: `postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable`
- `error`:  if any variable is missing or `""` empty string<br><br>


### Function: `HandlePgErrorFn(table string, err error) (int, error)`
Maps postgres error codes to http codes, table for example `user`<br>

Logic:
- Switch case that maps pg error code to http status code
- 23505 -> 409, 23514 -> 422, 42703 -> 400, 500 -> not mapped/unexpected<br>

Returns:
- `int`:    http status code
- `error`:  if execution wasn't successful + explantion why<br><br>


### Wrapper: `HandleSelectErrorFn(err error) (int, error)`
Check if query returned no rows or failed to execute.<br>

Logic:
- Check for `sql.ErrNoRows()`
- Check for errors<br>

Returns:
- `int`:    http status code
- `error`:  if query wasn't executed + reason why<br><br>


### Function: `CheckRowsAffectedInsertFn(result sql.Result) error`
Check if rows affected is different from exactly one !=1<br>

Returns:
- `error`:  if unexpcted number of rows affected<br><br>
<!-- }}} DB-->


## Models
<!-- {{{ Models -->
<!-- {{{ userModel -->
### Struct: `User`
struct for user with validate function.<br><br>


### Function: `IsValidUsernameFn(username string) error`
checks if username is correct lenght and dosen't have illegal chars<br>

Returns:
- `error`: if dosen't meet requirements + explanation why<br><br>


### Function: `isValidHexStringFn(hexStr string, hexStrName string, length int) error`
checks if string is correct lenght and matches HEX string.<br>

Returns:
- `error`: if dosen't meet requirements + explanation why<br><br>


### Wrapper: `.Validate() error`
checks if instance is valid.<br>

Requirements:
- instance: [`User`](shared.md#struct-user)
- function: [`IsValidUsernameFn()`](shared.md#function-isvalidusernamefnusername-string-error)
- function: [`isValidHexStringFn()`](shared.md#function-isvalidhexstringfnhexstr-string-hexstrname-string-length-int-error)<br>

Logic:
- Call `isValidusernameFn()`
- Call `isValidHexStringFn()` for `salt`, `hash`, `enc_symkey`<br>

Returns:
- `error`: if dosen't meet requirements + explanation why<br><br>
<!-- }}} userModel -->
<!-- }}} Models -->


## API
<!-- {{{ API -->
<!-- {{{ functions -->
### Function: `ExtractJSONValueFn(r *http.Request, key string, target interface{}) error`
Extracts value via key from JSON-encoded request body.<br>

Requirements:
- `encoding/json`<br>

Returns:
- `error`: if failed to parse request and extract value<br>

Side effects: if successful update target (with value)<br><br>
<!-- }}} functions -->


<!-- {{{ response -->
#TODO: update it
### Struct: `APIResponse`
struct for writing API response<br><br>


### Function: `WriteJSONResponseFn(w http.ResponseWrite, statusCode int, message string, errMsg string, data interface{})`
Writes uniformed API response based on give args

Requirements:
- `encoding/json`
- `net/http`<br>

Logic:
- Set Header: Content-Type: application/json
- Create `APIResponse` struct
- Encode it to JSON<br><br>
<!-- }}} response -->
<!-- }}} API -->


