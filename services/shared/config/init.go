package mconfig

import (
	"context"
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"google.golang.org/grpc"
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
func InitFromConfigCenter(serviceID string) error {
	// init service public config
	err := InitFromFile("../config.json")
	if err != nil {
		log.Println("init from file failed, error: ", err)
		return err
	}

	err = initServicePublicConfig()
	if err != nil {
		log.Println("init service public config failed, error: ", err)
		return err
	}

	// init rpc client for config center service and get config
	conn, err := grpc.Dial(GetCoreTarget(), grpc.WithInsecure())
	if err != nil {
		log.Println("get grpc conn failed, error: ", err)
		return err
	}

	var res *rpc_impl.ConfigCenter_GetServiceConfigRes
	for i := 0; i < servicePublicConfigIns.RetryTimes; i++ {
		res, err = rpc_impl.NewIConfigCenterClient(conn).GetServiceConfig(context.Background(), &rpc_impl.ConfigCenter_GetServiceConfigReq{
			ServiceId: serviceID,
			Level:     GetConfigLevel(),
		})
		if err == nil && (res != nil && res.Err == nil) { // no error, stop re-try
			break
		}

		log.Printf("get config from config center failed, retry in %d seconds in %d times, err: %v\n",
			servicePublicConfigIns.RetryInterval, servicePublicConfigIns.RetryTimes-i-1, err)

		time.Sleep(time.Duration(servicePublicConfigIns.RetryInterval) * time.Second)
	}
	if err != nil {
		log.Println("get config from config center failed, with grpc connection error: ", err)
		return err
	}
	if res == nil || res.Err != nil {
		log.Println("get config from config center failed, error: ", res.Err.String())
		return utils.NewError(res.Err.String())
	}

	err = json.Unmarshal([]byte(res.Config), conf)
	if err != nil {
		log.Println("json unmarshal failed, error: ", err)
		return err
	}

	return nil
}

// InitFromFile load config from config file, use for config center and test
func InitFromFile(configFileName string) error {
	if len(conf.Level) > 0 { // have initialized, require any config version have its level
		log.Println("already initialized")
		return nil
	}

	if len(configFileName) < 1 {
		configFileName = mconst.ConfigFileName
	}

	configBytes, err := os.ReadFile(configFileName)
	if err != nil {
		log.Printf("read config file failed, path: %s, error: %v\n", configFileName, err)
		return err
	}

	err = json.Unmarshal(configBytes, conf)
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
