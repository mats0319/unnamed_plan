package config

import (
	"encoding/json"
	"fmt"
	. "github.com/mats9693/unnamed_plan/shared/go/const"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	Level       string        `json:"level"`
	ConfigItems []*configItem `json:"config"`
}

type configItem struct {
	Name    string          `json:"name"`
	UID     string          `json:"uid"`
	Json    json.RawMessage `json:"json"`
	JsonDev json.RawMessage `json:"jsonDev"`
}

var conf = &config{}

func init() {
	str, err := os.Executable()
	if err != nil {
		fmt.Println("get executable failed, error:", err)
		return
	}

	dir := filepath.Dir(str)
	configPath := dir + "/" + ConfigFileName

	configBs, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("read config file failed, config path: %s, error:%v\n", configPath, err)
		os.Exit(-1)
	}

	err = json.Unmarshal(configBs, conf)
	if err != nil {
		fmt.Println("json unmarshal failed, error:", err)
		os.Exit(-1)
	}

	if conf.Level == ConfigDevLevel {
		for i := range conf.ConfigItems {
			conf.ConfigItems[i].Json = conf.ConfigItems[i].JsonDev
		}
	}

	fmt.Println("> Config init finish.")
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
