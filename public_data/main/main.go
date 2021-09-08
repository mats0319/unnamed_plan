package main

import (
	_ "github.com/mats9693/unnamed_plan/admin_data/config"
	_ "github.com/mats9693/unnamed_plan/admin_data/db"
	h2 "github.com/mats9693/unnamed_plan/admin_data/http"
)

func main() {
	h2.StartServer()
}
