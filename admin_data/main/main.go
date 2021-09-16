package main

import (
	_ "github.com/mats9693/unnamed_plan/admin_data/http"
	_ "github.com/mats9693/unnamed_plan/shared/go/config"
	_ "github.com/mats9693/unnamed_plan/shared/go/db"

	"github.com/mats9693/unnamed_plan/admin_data/config"
	"github.com/mats9693/unnamed_plan/shared/go/http"
)

func main() {
	err := system_config.InitConfiguration()
	if err != nil {
		return
	}

	shttp.StartServer()
}
