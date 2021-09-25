package system_config

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/utils/toy_server/config"
	. "github.com/mats9693/utils/toy_server/const"
)

type configuration struct {
	ARankAdminPermission uint8  `json:"ARankAdminPermission"`
	SRankAdminPermission uint8  `json:"SRankAdminPermission"`
	CloudFileRootPath    string `json:"cloudFileRootPath"`  // absolute path
	CloudFilePublicDir   string `json:"cloudFilePublicDir"` // public folder name
}

var systemConfig = &configuration{}

func GetConfiguration() *configuration {
	return systemConfig
}

func InitConfiguration() (err error) {
	byteSlice := config.GetConfig(UID_Config)

	err = json.Unmarshal(byteSlice, systemConfig)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", UID_Config, err)
		return
	}

	fmt.Println("> System config init finish.")

	return
}
