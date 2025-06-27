package middleware

import (
    "net/http"
    "fmt"
    "log"
)
import (
    sapi "github.com/FAH2S/diar4/src/shared/api"
)


func isMethodPOSTFn(r *http.Request) bool {
    return r.Method == http.MethodPost
}


func isHeaderCTAJFn(r *http.Request) bool {
    return r.Header.Get("Content-Type") == "application/json"
}


// Validate method and header wrapper
func ValidateMethodAndType(w http.ResponseWriter, r *http.Request) bool {
    const fn = "Middleware ValidateMethodAndType"
    ip := r.RemoteAddr
    if !isMethodPOSTFn(r){
        sapi.WriteJSONResponse(
            w,
            400,
            fmt.Sprintf("Fail: process '%s'", r.URL.Path),
            fmt.Sprintf("Method not allowed"),
            nil,
        )
        log.Printf("%s: Method not allowed | status: 400 | IP: %s", fn, ip)
        return false
    }
    // Check content type
    if !isHeaderCTAJFn(r){
        sapi.WriteJSONResponse(
            w,
            400,
            fmt.Sprintf("Fail: process '%s'", r.URL.Path),
            fmt.Sprintf("Content-Type must be application/json"),
            nil,
        )
        log.Printf("%s: Content-Type must be application/json | status: 400 | IP: %s", fn, ip)
        return false
    }
    return true
}


// Validate method and header endpoint/handler
func ValidateMethodAndTypeEndpoint(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !ValidateMethodAndType(w, r) {
            return
        }
        next.ServeHTTP(w, r)
    })
}
