package mconfig

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	Version     string        `json:"version"`
	Level       string        `json:"level"`
	ConfigItems []*configItem `json:"config"`
}

type configItem struct {
	UID  string          `json:"uid"`
	Name string          `json:"name"`
	Json json.RawMessage `json:"json"`
}

var (
	conf    = &config{}
	execDir = ""
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("get executable failed, error:", err)
		os.Exit(-1)
	}

	execDir = utils.FormatDirSuffix(dir)

	err = setConfig()
	if err != nil {
		os.Exit(-1)
	}

	log.Println("> Config init finish, more log will be redirect to log file.")
}

func setConfig() error {
	configPath := execDir + mconst.ConfigFileName

	configBs, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("read config file failed, path: %s, error: %v\n", configPath, err)
		return err
	}

	err = json.Unmarshal(configBs, conf)
	if err != nil {
		log.Println("json unmarshal failed, error:", err)
		return err
	}

	return nil
}

func GetConfigLevel() string {
	return conf.Level
}

func GetConfig(uid string) []byte {
	res := make([]byte, 0)
	for _, v := range conf.ConfigItems {
		if v.UID == uid {
			res = v.Json
			break
		}
	}

	return res
}

func GetExecDir() string {
	return execDir
}
