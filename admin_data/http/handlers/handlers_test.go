package handlers

import (
    "fmt"
    "os"
    "testing"
)

func TestMkdirAll(t *testing.T) {
    path := "D:/GoPath/src/github.com/mats9693/unnamed_plan/admin_data/main/a/b/c/"
    err := os.MkdirAll(path, 0755)

    fmt.Println("> --- Test MkdirAll ---")
    fmt.Println("> result:", existFile(path))
    fmt.Println("> error :", err)
}

func existFile(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}
