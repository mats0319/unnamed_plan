package config

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/utils/toy_server/config"
	"os"
)

const uid_CloudFileServiceConfig = "76175dee-ffee-4198-8f97-94dfd4f960fd"

type cloudFileServiceConfig struct {
	Address            string `json:"address"`
	UserServerAddress  string `json:"userServerAddress"`
	CloudFileRootPath  string `json:"cloudFileRootPath"`  // absolute path
	CloudFilePublicDir string `json:"cloudFilePublicDir"` // public folder name
}

var cloudFileServiceConfigIns = &cloudFileServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_CloudFileServiceConfig)

	err := json.Unmarshal(byteSlice, cloudFileServiceConfigIns)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", uid_CloudFileServiceConfig, err)
		os.Exit(-1)
	}

	fmt.Println("> User service config init finish.")
}

func GetConfig() *cloudFileServiceConfig {
	return cloudFileServiceConfigIns
}
