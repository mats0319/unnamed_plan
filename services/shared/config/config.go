package mconfig

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"log"
	"os"
)

type servicePublicConfig struct {
	init bool

	ConfigCenterTarget string `json:"configCenterTarget"`
	RetryTimes         int    `json:"retryTimes"`
	RetryInterval      int64  `json:"retryInterval"` // unit: second
}

var servicePublicConfigIns = &servicePublicConfig{}

func initServicePublicConfig() {
	if servicePublicConfigIns.init { // have initialized
		return
	}

	byteSlice := GetConfig(mconst.UID_Config)

	err := json.Unmarshal(byteSlice, servicePublicConfigIns)
	if err != nil {
		log.Printf("json unmarshal failed, uid: %s, error: %v\n", mconst.UID_Config, err)
		os.Exit(-1)
	}

	servicePublicConfigIns.init = true
}
