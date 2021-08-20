package main

import (
	_ "github.com/mats9693/unnamed_plan/public_data/config"
	_ "github.com/mats9693/unnamed_plan/public_data/db"
	"github.com/mats9693/unnamed_plan/public_data/http"
)

func main() {
	http.StartServer()
}