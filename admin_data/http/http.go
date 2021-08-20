package http

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/public_data/config"
	"net/http"
)

type httpConfig struct {
	Port string `json:"port"`
}

// StartServer is block
func StartServer() {
	conf := getHttpConfig()

	fmt.Printf("> Listening at : %s.", conf.Port)
	_ = http.ListenAndServe(":"+conf.Port, nil)
}

func getHttpConfig() *httpConfig {
	byteSlice := config.GetConfig("3b839c1f-9f1e-474b-b3da-7b5e9bc792ec")

	conf := &httpConfig{}
	err := json.Unmarshal(byteSlice, conf)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", "2", err)
		return nil
	}

	return conf
}
