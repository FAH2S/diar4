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
func CreateUserEndpoint(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var (
        wrap        = "CreateUserEndpoint"
        // Input
        user        smodels.User
        // Response info
        statusCode  = 500
        message     = "Fail: create user ''"
        errMessage  = "Unknown error occurred"
        ip          = r.RemoteAddr
        success     = false
    )

    // Helper fn, logs result and writes JSON response
    respond := func(err error) {
        if success {
            log.Printf("%s: %s | status: %d | IP: %s", wrap, message, statusCode, ip)
        } else {
            log.Printf("%s: %v | status: %d | IP: %s", wrap, err, statusCode, ip)
        }
        sapi.WriteJSONResponseFn(w, statusCode, message, errMessage, nil)
    }

    // Decode request body into user model
    err := json.NewDecoder(r.Body).Decode(&user); if err != nil {
        statusCode = 400
        errMessage = "Invalid JSON"
        respond(err); return
    }

    // Validate user model fields
    err = user.Validate(); if err != nil {
        statusCode = 422
        message = fmt.Sprintf("Fail: create user '%s'", user.Username)
        errMessage = fmt.Sprintf("Invalid input format: %v", err)
        respond(err); return
    }

    // Attempt to insert user
    statusCode, err = InsertUser(db, user)
    message, errMessage, success = sapi.MapStatusCodeFn(statusCode, "create", "user", user.Username, err)
    respond(err); return
}
//}}} Create user endpoint


//{{{ Read user endpoint
func ReadUserEndpoint(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var (
        wrap        = "ReadUserEndpoint"
        // Input
        username    = ""
        // Response info
        statusCode  = 500
        message     = "Fail: read user ''"
        errMessage  = "Unknown error occured"
        user        *smodels.User
        ip          = r.RemoteAddr
        success     = false
    )

    // Helper fn, logs result and writes JSON response
    respond := func(err error) {
        if success {
            log.Printf("%s: %s | status: %d | IP: %s", wrap, message, statusCode, ip)
        } else {
            log.Printf("%s: %v | status: %d | IP: %s", wrap, err, statusCode, ip)
        }
        sapi.WriteJSONResponseFn(w, statusCode, message, errMessage, user)
    }

    // Extract username from request body
    err := sapi.ExtractJSONValueFn(r, "username", &username); if err != nil {
        statusCode = 400
        errMessage = "Invalid JSON"
        respond(err); return
    }

    // Validate extracted username
    err = smodels.IsValidUsernameFn(username)
    if err != nil {
        statusCode = 422
        message = fmt.Sprintf("Fail: read user '%s'", username)
        errMessage = fmt.Sprintf("Invalid input format: %v", err)
        respond(err); return
    }

    // Attempt to select(fetch) user
    statusCode, user, err = SelectUser(db, username)
    message, errMessage, success = sapi.MapStatusCodeFn(statusCode, "read", "user", username, err)
    respond(err); return
}
//}}} Read user endpoint


//{{{ Update user endpoint
func UpdateUserEndpoint(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var (
        wrap        = "UpdateUserEndpoint"
        // Input
        inputData   map[string]interface{}
        username    = ""
        // Response info
        statusCode  = 500
        message     = "Fail: update user ''"
        errMessage  = "Unknown error occurred"
        returnData  map[string]string
        ip          = r.RemoteAddr
        success     = false
    )
    // Helper fn, logs result and writes JSON response
    respond := func(err error) {
        if success {
            log.Printf("%s: %s | status: %d | IP: %s", wrap, message, statusCode, ip)
        } else {
            log.Printf("%s: %v | status: %d | IP: %s", wrap, err, statusCode, ip)
        }
        sapi.WriteJSONResponseFn(w, statusCode, message, errMessage, returnData)
    }

    // Decode request body into map/dict data
    err := json.NewDecoder(r.Body).Decode(&inputData); if err != nil {
        statusCode = 400
        errMessage = "Invalid JSON"
        respond(err); return
    }

    // Sanitize data
    // - remove illegal keys => sapi
    allowed := []string{"username", "salt", "hash", "enc_symkey"}
    filterdData := sapi.SanitizeKeysFn(inputData, allowed)
    // - check each present field via some user validate => smodels
    err = smodels.ValidateUserMap(filterdData); if err != nil {
        statusCode = 422
        errMessage = fmt.Sprintf("Invalid input format: %v", err)
        respond(err); return
    }
    // - username field must be present, validate makes sure its string
    username, ok := filterdData["username"].(string)
    if !ok {
        statusCode = 422
        errMessage = "Missing required field: 'username'"
        respond(err); return
    }
    // - check for at least 2 fields
    if len(filterdData) < 2 {
        statusCode = 422
        errMessage = "Invalid input: must contain at least 2 fields total"
        respond(err); return
    }
    // - pop username from map
    delete (filterdData, "username")


    // Call UpdateUser
    statusCode, err = UpdateUser(db, filterdData, username)
    if statusCode == 200 {
        returnData = map[string]string{"username":username}
    }
    message, errMessage, success = sapi.MapStatusCodeFn(statusCode, "update", "user", username, err)
    // Return API response
    respond(err); return
}


//}}} Update user endpoint



