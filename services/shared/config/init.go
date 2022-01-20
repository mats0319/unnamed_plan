package mconfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
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
	conf = &config{}

	help          bool
	configFileDir string
	logFileDir    string
)

func init() {
	parseFlag()

	err := setDir()
	if err != nil {
		os.Exit(-1)
	}

	err = initConfig()
	if err != nil {
		os.Exit(-1)
	}

	log.Println("> Config init finish, more log will be redirect to log file.")
}

func parseFlag() {
	flag.BoolVar(&help, "h", false, "this help")
	flag.StringVar(&configFileDir, "conf", "", "config file dir")
	flag.StringVar(&logFileDir, "log", "", "log file dir")

	flag.Parse()

	if help {
		log.Println("Options: ")
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func setDir() error {
	if len(configFileDir) > 0 && len(logFileDir) > 0 {
		return nil
	}

	executableAbsolutePath, err := os.Executable()
	if err != nil {
		fmt.Println("get executable failed, error:", err)
		return err
	}
	executableAbsolutePath = strings.ReplaceAll(executableAbsolutePath, "\\", "/")
	executableDir := utils.FormatDirSuffix(path.Dir(executableAbsolutePath))

	if len(configFileDir) < 1 {
		configFileDir = executableDir
	}
	if len(logFileDir) < 1 {
		logFileDir = executableDir
	}

	return nil
}

func initConfig() error {
	configPath := configFileDir + mconst.ConfigFileName

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

func GetLogFileDir() string {
	return logFileDir
}
