package plugins

import (
	"encoding/json"
	"github.com/mats9693/unnamed_plan/services/gateway/plugins/limit_multi_login"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
)

func LoadValidPlugins() ([]mhttp.Plugins, error) {
	// get plugins config
	byteSlice := mconfig.GetConfig(mconst.UID_Gateway_Plugins)

	pluginsConfig := &struct {
		LimitMultiLoginConfig *limit_multi_login.LimitMultiLoginConfig `json:"limitMultiLoginConfig"`
	}{}

	err := json.Unmarshal(byteSlice, pluginsConfig)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_Gateway_Plugins), zap.Error(err))
		return nil, err
	}

	// load plugins which ones has config
	pluginsSlice := make([]mhttp.Plugins, 0)
	{
		if pluginsConfig.LimitMultiLoginConfig != nil {
			item := limit_multi_login.Init(pluginsConfig.LimitMultiLoginConfig)
			pluginsSlice = append(pluginsSlice, item)
		}
	}

	return pluginsSlice, nil
}
