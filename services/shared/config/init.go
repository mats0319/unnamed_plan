package mconfig

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"log"
	"os"
)

type config struct {
	Level       string        `json:"level"`
	ConfigItems []*configItem `json:"config"`
}

type configItem struct {
	UID  string          `json:"uid"`
	Name string          `json:"name"`
	Json json.RawMessage `json:"json"`
}

var conf = &config{}

// InitFromConfigCenter load config from config center, use for services
func InitFromConfigCenter(serviceID string, level string) {
	// todo: impl
}

// InitFromFile load config from config file, use for config center and test
func InitFromFile(configFileName string) {
	if len(conf.Level) > 0 { // have initialized, require any config version have its level
		return
	}

	if len(configFileName) < 1 {
		configFileName = mconst.ConfigFileName
	}

	configBytes, err := os.ReadFile(configFileName)
	if err != nil {
		log.Printf("read config file failed, path: %s, error: %v\n", configFileName, err)
		os.Exit(-1)
	}

	err = json.Unmarshal(configBytes, conf)
	if err != nil {
		log.Println("json unmarshal failed, error:", err)
		os.Exit(-1)
	}

	log.Println("> Config init finish, more log will be redirect to log file.")
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
