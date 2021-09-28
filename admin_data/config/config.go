package system_config

import (
	"encoding/json"
	"fmt"
	"os"

	mconfig "github.com/mats9693/utils/toy_server/config"
	mconst "github.com/mats9693/utils/toy_server/const"
)

type configuration struct {
	ARankAdminPermission uint8  `json:"ARankAdminPermission"`
	SRankAdminPermission uint8  `json:"SRankAdminPermission"`
	CloudFileRootPath    string `json:"cloudFileRootPath"`  // absolute path
	CloudFilePublicDir   string `json:"cloudFilePublicDir"` // public folder name
}

var systemConfig = &configuration{}

func init() {
	byteSlice := mconfig.GetConfig(mconst.UID_Config)

	err := json.Unmarshal(byteSlice, systemConfig)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", mconst.UID_Config, err)
		os.Exit(-1)
	}

	fmt.Println("> System config init finish.")
}

func GetConfiguration() *configuration {
	return systemConfig
}
