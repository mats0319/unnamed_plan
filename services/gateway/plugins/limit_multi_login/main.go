package limit_multi_login

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"net/http"
	"sync"
)

type LimitMultiLoginConfig struct {
	LimitMultiLogin bool     `json:"limitMultiLogin"`
	Sources         []string `json:"sources"`
	KeepTokenValid  int64    `json:"KeepTokenValid"` // unit: second
}

type LimitMultiLogin struct {
	config *LimitMultiLoginConfig

	loginInfoMap sync.Map // user id - login info

	flags sync.Map // pattern - flags
}

type flags struct {
	skipMultiLoginLimit   bool
	reSetMultiLoginParams bool
}

var limitMultiLoginIns = &LimitMultiLogin{}

func Init(config *LimitMultiLoginConfig) mhttp.Plugins {
	limitMultiLoginIns.config = config

	return limitMultiLoginIns
}

func HandleFunc(pattern string, params ...uint8) {
	flagsIns := &flags{}

	for i := range params {
		switch params[i] {
		case mconst.HTTPFlags_MultiLogin_SkipLimit:
			flagsIns.skipMultiLoginLimit = true
		case mconst.HTTPFlags_MultiLogin_ReSetParams:
			flagsIns.reSetMultiLoginParams = true
		}
	}

	limitMultiLoginIns.flags.Store(pattern, flagsIns)
}

func (l *LimitMultiLogin) RunBeforeHook(uri string) bool {
	return l.runHook(uri, true)
}

func (l *LimitMultiLogin) RunAfterHook(uri string) bool {
	return l.runHook(uri, false)
}

func (l *LimitMultiLogin) BeforeHook(request *http.Request, timestamp int64) (int, error, bool) {
	source := request.Header.Get(mconst.HTTP_SourceSign)
	if !utils.Contains(l.config.Sources, source) { // ignore unlimited sources
		return http.StatusOK, nil, false
	}

	userID := request.Header.Get(mconst.HTTP_MultiLogin_UserIDReq)
	token := request.Header.Get(mconst.HTTP_MultiLogin_TokenReq)

	err := l.multiLoginVerify(userID, source, token, timestamp)
	if err != nil {
		return http.StatusUnauthorized, err, true
	}

	return http.StatusOK, nil, false
}

func (l *LimitMultiLogin) AfterHook(writer http.ResponseWriter, request *http.Request, timestamp int64, param string) (int, error, bool) {
	source := request.Header.Get(mconst.HTTP_SourceSign)
	if !utils.Contains(l.config.Sources, source) { // ignore unlimited sources
		return http.StatusOK, nil, false
	}

	newToken := utils.RandomHexString(10)

	err := l.setLoginInfo(param, source, newToken, timestamp)
	if err != nil {
		return http.StatusInternalServerError, err, true
	}

	writer.Header().Set(mconst.HTTP_MultiLogin_UserIDRes, param)
	writer.Header().Set(mconst.HTTP_MultiLogin_TokenRes, newToken)

	return http.StatusOK, nil, false
}
