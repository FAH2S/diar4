

## DB
<!-- {{{ DB -->
### Wrapper: `db.GetConn() (*sql.DB, error)`
Creates new connection to database from environment variables.

Requirements:
- `DB_USER, DB_PWD, DB_HOST, DB_PORT, DB_NAME` must be set as non empty string.
- function: [`buildConnStrFromEnv()`](#function-buildconnstrfromenvfn)
- sql.Open

Logic:
- Call buildCononStrFromEnvFn() to create connection string
- Call sql.Open to create DB connection
- Return connection or error

Returns:
- `*sql.db`:    pointer to db connection
- `error`:      if any variable is missing or `sql.open` fails


### Function: `buildConnStrFromEnvFn()`
Builds connection string from environment variables.

Logic:
- Read enviroment variables + validate
- Format connection string from enviroment variables

Returns:
- `string`: connection string ex.:
- `error`:  if any variable is missing or `""` empty string


### Function: `HandlePgError(err error) (int, error)`
Maps postgres error codes to http codes

Logic:
- Switch case that maps pg error code to http status code
- 23505 -> 409, 23514 -> 422, 42703 -> 400, 500 -> not mapped/unexpected

Returns:
- `int`:    http status code
- `error`:  if execution wasn't successful + explantion why


### Function: `CheckRowsAffectedInsert(rows int64) error`
Check if rows affected is not zero

Returns:
- `error`:  if unexpcted number of rows affected
<!-- }}} -->


## Models
<!-- {{{ Models -->
<!-- {{{ userModel -->
### Struct: `userModel`
struct for user with validate function.


### Function: `isValidUsernameFn(username string) error`
checks if username is correct lenght and dosen't have illegal chars

Returns:
- `error`: if dosen't meet requirements + explanation why


### Function: `isValidHexStringFn(hexStr string, hexStrName string, length int) error`
checks if string is correct lenght and matches HEX string.

Returns:
- `error`: if dosen't meet requirements + explanation why


### Wrapper: `.validate() error`
checks if instance is valid.

Requirements:
- instance: [`userModel`](#struct-usermodel)
- function: [`isValidUsernameFn()`](#function-isvalidusernamefnusername-stringerror)
- function: [`isValidHexStringFn()`](#function-isvalidhexstringfnhexstr-string-hexstrname-string-length-interror)

Returns:
- `error`: if dosen't meet requirements + explanation why
<!-- }}} userModel -->
<!-- }}} Models -->

