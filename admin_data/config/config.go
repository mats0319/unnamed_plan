package system_config

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/shared/go/config"
	. "github.com/mats9693/unnamed_plan/shared/go/const"
)

type configuration struct {
	ARankAdminPermission uint8 `json:"ARankAdminPermission"`
	SRankAdminPermission uint8 `json:"SRankAdminPermission"`
}

var systemConfig = &configuration{}

func GetConfiguration() *configuration {
	// todo: consider if use singleton pattern
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
