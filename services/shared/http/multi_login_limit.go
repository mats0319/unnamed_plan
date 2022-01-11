package mhttp

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"net/http"
)

type loginInfo struct {
	tokens     []string
	timestamps []int64
}

// multiLoginVerify verify login info of designated user, if passed, re-set timestamp
func (h *Handlers) multiLoginVerify(userID string, source string, token string, timestamp int64) string {
	if len(userID) < 1 || len(token) < 1 || timestamp <= 0 {
		return fmt.Sprintf("%s. user id: %s, token: %s, timestamp: %d",
			mconst.Error_InvalidParams, userID, token, timestamp)
	}

	// get login info of all sources by 'userID'
	loginInfoIns := &loginInfo{}
	{
		v, ok := h.loginInfoMap.Load(userID)
		if !ok || len(userID) < 1 {
			return fmt.Sprintf("%s, user id: %s", mconst.Error_LoadLoginInfoFailed, userID)
		}

		loginInfoIns, ok = v.(*loginInfo)
		if !ok || len(loginInfoIns.tokens) < 1 || len(loginInfoIns.timestamps) < 1 {
			return http.StatusText(http.StatusInternalServerError)
		}
	}

	// get login info of designated source by 'source'
	index := utils.GetIndex(h.config.Sources, source)
	if index < 0 {
		return mconst.Error_UnknownSource + source
	}

	cachedToken := loginInfoIns.tokens[index]
	cachedTimestamp := loginInfoIns.timestamps[index]

	// do multi-login limit verify
	errMsg := ""
	if token != cachedToken {
		errMsg += fmt.Sprintf(" { %s: %s } ", mconst.Error_InvalidToken, token)
	}
	if timestamp > cachedTimestamp+h.config.KeepTokenValid {
		errMsg += fmt.Sprintf(" { %s } ", mconst.Error_InvalidTokenTimeout)
	}

	if len(errMsg) > 0 {
		return errMsg
	}

	// re-set 'timestamp'
	loginInfoIns.timestamps[index] = timestamp
	h.loginInfoMap.Store(userID, loginInfoIns)

	return ""
}

// setLoginInfo set login info, only set valid params(not zero value)
func (h *Handlers) setLoginInfo(userID string, source string, token string, timestamp int64) string {
	if len(userID) < 1 || (len(token) < 1 && timestamp <= 0) {
		return fmt.Sprintf("%s. user id: %s, token: %s, timestamp: %d",
			mconst.Error_InvalidParams, userID, token, timestamp)
	}

	index := utils.GetIndex(h.config.Sources, source)
	if index < 0 {
		return mconst.Error_UnknownSource + source
	}

	// get login info of all sources by 'userID' or init if not exist
	loginInfoIns := &loginInfo{}
	if v, ok := h.loginInfoMap.Load(userID); ok {
		loginInfoIns, _ = v.(*loginInfo)
	} else {
		loginInfoIns.tokens = make([]string, len(h.config.Sources))
		loginInfoIns.timestamps = make([]int64, len(h.config.Sources))
	}

	// set valid login info(not zero value)
	if len(token) > 0 {
		loginInfoIns.tokens[index] = token
	}
	if timestamp > 0 {
		loginInfoIns.timestamps[index] = timestamp
	}

	h.loginInfoMap.Store(userID, loginInfoIns)

	return ""
}
