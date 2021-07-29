package main

import (
	_ "github.com/mats9693/unnamed_plan/config"
	_ "github.com/mats9693/unnamed_plan/db"
	"github.com/mats9693/unnamed_plan/http"
)

func main() {
	http.StartServer()
}
