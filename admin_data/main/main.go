package main

import (
	_ "github.com/mats9693/unnamed_plan/admin_data/http"
	_ "github.com/mats9693/unnamed_plan/shared/go/config"
	_ "github.com/mats9693/unnamed_plan/shared/go/db"
	"github.com/mats9693/unnamed_plan/shared/go/http"
)

func main() {
	shttp.StartServer()
}
