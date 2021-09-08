package http

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/config"
	. "github.com/mats9693/unnamed_plan/admin_data/const"
	"net/http"
)

type httpConfig struct {
	Port string `json:"port"`
}

// StartServer is block
func StartServer() {
	conf := getHttpConfig()

	fmt.Printf("> Listening at : %s.\n", conf.Port)
	_ = http.ListenAndServe(":"+conf.Port, nil)
}

func getHttpConfig() *httpConfig {
	byteSlice := config.GetConfig(UID_HTTPPort)

	conf := &httpConfig{}
	err := json.Unmarshal(byteSlice, conf)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", UID_HTTPPort, err)
		return nil
	}

	return conf
}
