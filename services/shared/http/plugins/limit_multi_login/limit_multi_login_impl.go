package limit_multi_login

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/http/plugins"
	"github.com/mats9693/unnamed_plan/services/shared/http/response"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"net/http"
	"sync"
)

type LimitMultiLoginConfig struct {
	LimitMultiLogin bool  `json:"limitMultiLogin"`
	KeepTokenValid  int64 `json:"KeepTokenValid"` // unit: second
}

type LimitMultiLogin struct {
	config  *LimitMultiLoginConfig
	sources []string

	loginInfoMap sync.Map // user id - login info

	flags sync.Map // pattern - flags
}

type flags struct {
	skipMultiLoginLimit   bool
	reSetMultiLoginParams bool
}

var limitMultiLoginIns = &LimitMultiLogin{}

func Init(config *LimitMultiLoginConfig, sources []string) plugins.Plugins {
	if limitMultiLoginIns.config == nil {
		limitMultiLoginIns.config = config
		limitMultiLoginIns.sources = sources
	}

	return limitMultiLoginIns
}

func HandleFunc(pattern string, params ...uint8) {
	flagsIns := &flags{}

	for i := range params {
		switch params[i] {
		case mconst.HTTPMultiLogin_SkipLimit:
			flagsIns.skipMultiLoginLimit = true
		case mconst.HTTPMultiLogin_ReSetParams:
			flagsIns.reSetMultiLoginParams = true
		}
	}

	limitMultiLoginIns.flags.Store(pattern, flagsIns)
}

func (l *LimitMultiLogin) BeforeInvokeHook(request *http.Request, timestamp int64) (int, error) {
	source := request.Header.Get(mconst.HTTP_SourceSign)
	if !utils.Contains(l.sources, source) { // ignore unlimited sources
		return http.StatusOK, nil
	}

	flagsIns, code, err := l.getFlags(request.RequestURI)
	if err != nil {
		return code, err
	}

	if flagsIns.skipMultiLoginLimit {
		return http.StatusOK, nil
	}

	userID := request.Header.Get(mconst.HTTP_MultiLoginUserIDSign)
	token := request.Header.Get(mconst.HTTP_MultiLoginTokenSign)

	err = l.multiLoginVerify(userID, source, token, timestamp)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	return http.StatusOK, nil
}

func (l *LimitMultiLogin) AfterInvokeHook(response *mresponse.ResponseData, request *http.Request, timestamp int64) (int, error) {
	source := request.Header.Get(mconst.HTTP_SourceSign)
	if !utils.Contains(l.sources, source) { // ignore unlimited sources
		return http.StatusOK, nil
	}

	flagsIns, code, err := l.getFlags(request.RequestURI)
	if err != nil {
		return code, err
	}

	if !flagsIns.reSetMultiLoginParams {
		return http.StatusOK, nil
	}

	newToken := utils.RandomHexString(10)

	err = l.setLoginInfo(response.UserID, source, newToken, timestamp)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	response.Token = newToken

	return http.StatusOK, nil
}
