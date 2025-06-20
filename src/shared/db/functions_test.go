package shareddb
import (
    "os"
    "testing"
)

//{{{ buildConnStrFromEnv
// succ
func TestBuildConnStrFromEnvSucc(t *testing.T){
    os.Clearenv()
    // Set required env vars for test
    os.Setenv("DB_USER", "testuser")
    os.Setenv("DB_PWD", "testpass")
    os.Setenv("DB_NAME", "testdb")
    os.Setenv("DB_HOST", "localhost")
    os.Setenv("DB_PORT", "5432")

    // expected outcome
    expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
    // check
    actual, err := buildConnStrFromEnv()
    if err != nil {
        t.Fatalf("Fatal, expected no error, got: %v", err)
    }
    if actual != expected {
        t.Errorf("Wrong return value:\nExpected: %s\nGot: %s", expected, actual)
    }
}


// fail, missing env
func TestBuildConnStrFromEnvMissingEnv(t *testing.T){
    // Clear env to cause err
    os.Clearenv()

    // expected outcome
    expected := ""
    expectedErr := "buildConnStrFromEnv: Required environment variable DB_USER, not set\n"
    // check
    actual, err := buildConnStrFromEnv()
    if err == nil {
        t.Fatalf("Fatal, epected error, got nil")
    }
    if actual != expected {
        t.Errorf("Wrong string value:\nExpected: %s\nGot: %s", expected, actual)
    }
    if err.Error() != expectedErr {
        t.Errorf("Wrong error value:\nExpected: %s\nGot: %s", expectedErr, err.Error())
    }
}
//}}} buildConnStrFromEnv


//{{{ HandlePgError
// Don't see point in unit testing
//}}} HandlePgError


//{{{ CheckRowsAffectedInsert
// Don't see point in unit testing
//}}} CheckRowsAffectedInsert


