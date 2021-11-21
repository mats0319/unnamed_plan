package config

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/utils/toy_server/config"
	"os"
)

const uid_ThinkingNoteServiceConfig = "5d71f9af-63f5-44c5-9ef0-403bcd1ad381"

type thinkingNoteServiceConfig struct {
	Address           string `json:"address"`
	UserServerAddress string `json:"userServerAddress"`
}

var thinkingNoteServiceConfigIns = &thinkingNoteServiceConfig{}

func init() {
	byteSlice := mconfig.GetConfig(uid_ThinkingNoteServiceConfig)

	err := json.Unmarshal(byteSlice, thinkingNoteServiceConfigIns)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", uid_ThinkingNoteServiceConfig, err)
		os.Exit(-1)
	}

	fmt.Println("> User service config init finish.")
}

func GetConfig() *thinkingNoteServiceConfig {
	return thinkingNoteServiceConfigIns
}
