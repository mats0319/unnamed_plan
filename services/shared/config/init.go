package mconfig

import (
	"context"
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/proto/client"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"log"
	"os"
	"time"
)

type Config struct {
	Level       string        `json:"level"`
	ConfigItems []*ConfigItem `json:"config"`
}

type ConfigItem struct {
	UID  string          `json:"uid"`
	Name string          `json:"name"`
	Json json.RawMessage `json:"json"`
}

var conf = &Config{}

// InitFromConfigCenter load config from config center, use for services
func InitFromConfigCenter(serviceID string) {
	// init service public config
	InitFromFile("../config.json")

	initServicePublicConfig()

	// init rpc client for config center service and get config
	configCenterClient, err := client.ConnectConfigCenterServer(servicePublicConfigIns.ConfigCenterTarget)
	if err != nil {
		log.Printf("establish connection with services failed, error: %v", err)
		os.Exit(-1)
	}

	var res *rpc_impl.ConfigCenter_GetServiceConfigRes
	for i := 0; i < servicePublicConfigIns.RetryTimes; i++ {
		res, err = configCenterClient.GetServiceConfig(context.Background(), &rpc_impl.ConfigCenter_GetServiceConfigReq{
			ServiceId: serviceID,
			Level:     GetConfigLevel(),
		})
		if err == nil {
			break
		}

		if err != nil {
			log.Printf("get config from config center failed, retry in %d seconds in %d times, err: %v\n",
				servicePublicConfigIns.RetryInterval, servicePublicConfigIns.RetryTimes-i, err)
			time.Sleep(time.Duration(servicePublicConfigIns.RetryInterval) * time.Second)
		}
	}
	if err != nil || res == nil {
		log.Printf("get config from config center failed, error: %v", err)
		os.Exit(-1)
	}

	err = json.Unmarshal([]byte(res.Config), conf)
	if err != nil {
		log.Println("json unmarshal failed, error:", err)
		os.Exit(-1)
	}
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
