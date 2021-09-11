package config

import (
    "encoding/json"
    "fmt"
    . "github.com/mats9693/unnamed_plan/shared/go/const"
)

type Configuration struct {
    CreateUserPermission uint8 `json:"createUserPermission"`
}

var configuration = &Configuration{}

func GetConfiguration() *Configuration {
    return configuration
}

func initConfiguration() *Configuration {
    byteSlice := GetConfig(UID_Config)

    err := json.Unmarshal(byteSlice, configuration)
    if err != nil {
        fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", UID_Config, err)
        return nil
    }

    return configuration
}
