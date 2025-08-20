package sharedapi
import (
    "encoding/json"
    "net/http"
    "fmt"
    "strings"
)

//TODO: rethink about this fn current problem is reading value will dump/erase
//  r.Body either repouplate it or rethink if this is even needed
func ExtractJSONValueFn(r *http.Request, key string, target interface{}) error {
    fn := "ExtractJSONValueFn"
    var raw map[string]json.RawMessage
    if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
        return fmt.Errorf("%s: Invalid JSON: %w", fn, err)
    }

    val, ok := raw[key]
    if !ok {
        return fmt.Errorf("%s: missing key: %s", fn, key)
    }

    if err := json.Unmarshal(val, target); err != nil {
        return fmt.Errorf("%s: failed to unmarshal value: %v to target type: %T", fn, val, target)
    }

    return nil
}

//TODO: potentially not api-layer but something else
func SanitizeKeysFn(inputMap map[string]interface{}, allowed []string) (map[string]interface{}) {
    allowedSet := make(map[string]struct{}, len(allowed))
    // Poppulate allowedSet with allowed keys
    for _, key := range allowed {
        allowedSet[key] = struct{}{} // Use struct so ok isn't false
    }
    // Iterate over inputMap and check if its key is present in allowedSet
    for k := range inputMap {
        if _, ok := allowedSet[k]; !ok {
            delete(inputMap, k)
        }
    }
    return inputMap
}



func MapStatusCodeFn(code int, action, entity, name string, err error) (msg, errMsg string, isSucc bool) {
    switch code {
    case 200:
        msg := fmt.Sprintf("Success: %s %s '%s'", action, entity, name)
        return msg, "", true
    case 201:
        msg := fmt.Sprintf("Success: %s %s '%s'", action, entity, name)
        return msg, "", true
    case 404:
        msg := fmt.Sprintf("Fail: %s %s '%s'", action, entity, name)
        errMsg := fmt.Sprintf("%s not found, dosen't exist", strings.ToUpper(entity[:1]) + entity[1:])
        return msg, errMsg, false
    case 409:
        msg := fmt.Sprintf("Fail: %s %s '%s'", action, entity, name)
        errMsg := fmt.Sprintf("%s already exist", strings.ToUpper(entity[:1]) + entity[1:])
        return msg, errMsg, false
    case 422:
        msg := fmt.Sprintf("Fail: %s %s '%s'", action, entity, name)
        errMsg := fmt.Sprintf("Invalid input format: %v", err)
        return msg, errMsg, false
    case 500:
        msg := fmt.Sprintf("Fail: %s %s", action, entity)
        errMsg := "Internal server error"
        return msg, errMsg, false
    default:
        msg := fmt.Sprintf("Fail: %s %s", action, entity)
        errMsg := "Unknown error"
        return msg, errMsg, false
    }
}
