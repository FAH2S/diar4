package sharedapi
import (
    "encoding/json"
    "net/http"
    "fmt"
)


func ExtractJSONValue(r *http.Request, key string, target interface{}) error {
    fn := "ExtractJSONValue"
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
