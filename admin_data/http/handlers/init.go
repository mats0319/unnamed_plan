package handlers

import (
    mconfig "github.com/mats9693/utils/toy_server/config"
    mconst "github.com/mats9693/utils/toy_server/const"
)

var isDev bool

func init() {
    isDev = mconfig.GetConfigLevel() == mconst.ConfigDevLevel
}
