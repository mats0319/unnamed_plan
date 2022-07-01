package plugins

import (
	"github.com/mats9693/unnamed_plan/services/gateway/plugins/limit_multi_login"
	"github.com/mats9693/unnamed_plan/services/shared/http"
)

func LoadValidPlugins() []mhttp.Plugins {
	plugins := make([]mhttp.Plugins, 0)

	p, err := limit_multi_login.GetPlugins()
	if err == nil {
		plugins = append(plugins, p)
	}

	return plugins
}
