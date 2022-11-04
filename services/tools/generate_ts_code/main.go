package main

import (
    "fmt"
    "os"
)

func main() {
    dir, err := os.Getwd() // go.mod dir
    if err != nil {
        fmt.Println("get dir failed", err)
        return
    }
    //fmt.Println("> Node: show dir.", dir)

    fileBytes, err := os.ReadFile(dir + "/tools/generate_ts_code/proto/common.proto")
    if err != nil {
        fmt.Println("read file fail", err)
        return
    }

    fmt.Println(string(fileBytes))
}
