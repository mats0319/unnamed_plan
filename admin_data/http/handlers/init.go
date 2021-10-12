package handlers

import (
    "github.com/mats9693/utils/toy_server/config"
    "github.com/mats9693/utils/toy_server/const"
)

var isDev bool

func init() {
    isDev = mconfig.GetConfigLevel() == mconst.ConfigDevLevel
}