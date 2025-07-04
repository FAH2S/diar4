package cruduser
import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "database/sql"
)
import (
    smodels "github.com/FAH2S/diar4/src/shared/models"
    sapi "github.com/FAH2S/diar4/src/shared/api"
)


//{{{ Create user endpoint
// TODO: if you want it to be more WRAPPER like extract parts like:
//  decode user
//  validate
//  insert user
//  handle statusCode
//  log results (once)
//  write response aka return (once)
func CreateUserEndpoint(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var (
        fn          = "CreateUserEndpoint"
        user        = smodels.User{}
        statusCode  = 500
        message     = "Fail: create user ''"
        errMessage  = "Unknown error occured"
        ip          = r.RemoteAddr
    )
    // Extract packet as user type/model
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        log.Printf("%s: %v | status: %d | IP: %s", fn, err, 400, ip)
        errMessage = "Invalid JSON"
        sapi.WriteJSONResponse(w, 400, message, errMessage, nil)
        return
    }
    // Validate extracted user
    err = user.Validate()
    if err != nil {
        log.Printf("%s: user.Validate: %v | status: %d | IP: %s", fn, err, 422, ip)
        message = fmt.Sprintf("Fail: create user '%s'", user.Username)
        errMessage = fmt.Sprintf("Invalid input format: %v", err)
        sapi.WriteJSONResponse(w, 422, message, errMessage, nil)
        return
    }
    // Insert user
    statusCode, err = InsertUser(db, user)
    switch statusCode {
    case 201:
        message = fmt.Sprintf("Success: create user '%s'", user.Username)
        errMessage = ""
    case 409:
        message = fmt.Sprintf("Fail: create user '%s'", user.Username)
        errMessage = "User already exist"
    case 422:
        message = fmt.Sprintf("Fail: create user '%s'", user.Username)
        errMessage = fmt.Sprintf("Invalid input format: %v", err)
    case 500:
        errMessage = "Internal server error"
    default:
        log.Printf("%s: Unknown error occured: %v", fn, err)
    }
    if statusCode == 201 {
        log.Printf("%s: %s | status: %d | IP: %s", fn, message, statusCode, ip)
    } else {
        log.Printf("%s: %v | status: %d | IP: %s", fn, err, statusCode, ip)
    }
    sapi.WriteJSONResponse(w, statusCode, message, errMessage, nil)
    return
}
//}}} Create user endpoint


//{{{ Read user endpoint
func ReadUserEndpoint(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var (
        fn          = "ReadUserEndpoint"
        username    = ""
        statusCode  = 500
        message     = "Fail: read user ''"
        errMessage  = "Unknown error occured"
        ip          = r.RemoteAddr
    )

    // Extract username
    err := sapi.ExtractJSONValue(r, "username", &username)
    if err != nil {
        log.Printf("%s: %v | status: %d | IP: %s", fn, err, 400, ip)
        errMessage = "Invalid JSON"
        sapi.WriteJSONResponse(w, 400, message, errMessage, nil)
        return
    }

    // Is valid username
    err = smodels.IsValidUsernameFn(username)
    if err != nil {
        log.Printf("%s: user.Validate: %v | status: %d | IP: %s", fn, err, 422, ip)
        message = fmt.Sprintf("Fail: read user '%s'", username)
        errMessage = fmt.Sprintf("Invalid input format: %v", err)
        sapi.WriteJSONResponse(w, 422, message, errMessage, nil)
        return
    }

    // SelectUser
    statusCode, user, err := SelectUser(db, username)
    switch statusCode {
    case 200:
        message = fmt.Sprintf("Success: read user '%s'", username)
        errMessage = ""
        // TODO: log
    case 404:
        message = fmt.Sprintf("Fail: read user '%s'", username)
        errMessage = "User not found, dosen't exist"
        // TODO: log
    case 500:
        errMessage = "Internal server error"
        // TODO: log
    default:
        log.Printf("%s: Unknown error occured: %v", fn, err)
    }
    // return
    sapi.WriteJSONResponse(w, statusCode, message, errMessage, user)
    return
}
//}}} Read user endpoint


