package config

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/utils/toy_server/config"
	"os"
)

const uid_UserServiceConfig = "6d31fac1-346c-4b03-8596-8b3e1ee5b960"

type userServiceConfig struct {
	Address              string `json:"address"`
	ARankAdminPermission uint8  `json:"ARankAdminPermission"`
	SRankAdminPermission uint8  `json:"SRankAdminPermission"`
}

var userServiceConfigIns = &userServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_UserServiceConfig)

	err := json.Unmarshal(byteSlice, userServiceConfigIns)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", uid_UserServiceConfig, err)
		os.Exit(-1)
	}

	fmt.Println("> User service config init finish.")
}

func GetConfig() *userServiceConfig {
	return userServiceConfigIns
}
