package sharedapi
import (
    "net/http"
    "encoding/json"
)


type APIResponse struct {
    Message string      `json:"message"`
    Error   string      `json:"error"`
    Data    interface{} `json:"data"`
}


func WriteJSONResponseFn(w http.ResponseWriter, statusCode int, message string, errMsg string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

    response := APIResponse{
        Message:    message,
        Error:      errMsg,
        Data:       data,
    }
    json.NewEncoder(w).Encode(response)
}
