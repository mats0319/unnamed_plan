package main

import (
	_ "github.com/mats9693/unnamed_plan/admin_data/config"
	_ "github.com/mats9693/unnamed_plan/admin_data/http"
	_ "github.com/mats9693/utils/toy_server/config"
	_ "github.com/mats9693/utils/toy_server/db"

	mhttp "github.com/mats9693/utils/toy_server/http"
)

func main() {
	mhttp.StartServer()
}
