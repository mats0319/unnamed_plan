package limit_multi_login

import (
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/pkg/errors"
	"net/http"
)

type loginInfo struct {
	tokens     []string
	timestamps []int64
}

// multiLoginVerify verify login info of designated user, if passed, re-set timestamp
func (l *LimitMultiLogin) multiLoginVerify(userID string, source string, token string, timestamp int64) error {
	if len(userID) < 1 || len(token) < 1 || timestamp <= 0 {
		err := errors.New(fmt.Sprintf("%s. user id: %s, token: %s, timestamp: %d",
			mconst.Error_InvalidParams, userID, token, timestamp))
		mlog.Logger().Error(err.Error())
		return err
	}

	// get login info of all sources by 'userID'
	loginInfoIns := &loginInfo{}
	{
		v, ok := l.loginInfoMap.Load(userID)
		if !ok || len(userID) < 1 {
			err := errors.New(fmt.Sprintf("%s, user id: %s", mconst.Error_LoadLoginInfoFailed, userID))
			mlog.Logger().Warn(err.Error())
			return err
		}

		loginInfoIns, ok = v.(*loginInfo)
		if !ok || len(loginInfoIns.tokens) < 1 || len(loginInfoIns.timestamps) < 1 {
			err := errors.New(http.StatusText(http.StatusInternalServerError))
			mlog.Logger().Error(err.Error())
			return err
		}
	}

	// get login info of designated source by 'source'
	index := utils.GetIndex(l.config.Sources, source)
	if index < 0 {
		err := errors.New(mconst.Error_UnknownSource + source)
		mlog.Logger().Error(err.Error())
		return err
	}

	cachedToken := loginInfoIns.tokens[index]
	cachedTimestamp := loginInfoIns.timestamps[index]

	// do multi-login limit verify
	errMsg := ""
	if token != cachedToken {
		errMsg += fmt.Sprintf(" { %s: %s } ", mconst.Error_InvalidToken, token)
	}
	if timestamp > cachedTimestamp+l.config.KeepTokenValid {
		errMsg += fmt.Sprintf(" { %s } ", mconst.Error_InvalidTokenTimeout)
	}

	if len(errMsg) > 0 {
		mlog.Logger().Warn(errMsg)
		return errors.New(errMsg)
	}

	// re-set 'timestamp'
	loginInfoIns.timestamps[index] = timestamp
	l.loginInfoMap.Store(userID, loginInfoIns)

	return nil
}

// setLoginInfo set login info, only set valid params(not zero value)
func (l *LimitMultiLogin) setLoginInfo(userID string, source string, token string, timestamp int64) error {
	if len(userID) < 1 || (len(token) < 1 && timestamp <= 0) {
		err := errors.New(fmt.Sprintf("%s. user id: %s, token: %s, timestamp: %d",
			mconst.Error_InvalidParams, userID, token, timestamp))
		mlog.Logger().Error(err.Error())
		return err
	}

	index := utils.GetIndex(l.config.Sources, source)
	if index < 0 {
		err := errors.New(mconst.Error_UnknownSource + source)
		mlog.Logger().Error(err.Error())
		return err
	}

	// get login info of all sources by 'userID' or init if not exist
	loginInfoIns := &loginInfo{}
	if v, ok := l.loginInfoMap.Load(userID); ok {
		loginInfoIns, _ = v.(*loginInfo)
	} else {
		loginInfoIns.tokens = make([]string, len(l.config.Sources))
		loginInfoIns.timestamps = make([]int64, len(l.config.Sources))
	}

	// set valid login info(not zero value)
	if len(token) > 0 {
		loginInfoIns.tokens[index] = token
	}
	if timestamp > 0 {
		loginInfoIns.timestamps[index] = timestamp
	}

	l.loginInfoMap.Store(userID, loginInfoIns)

	return nil
}

func (l *LimitMultiLogin) getFlags(pattern string) *flags {
	flagsI, ok := l.flags.Load(pattern)
	if !ok {
		mlog.Logger().Info(mconst.Error_UnknownURI + pattern)
		return nil
	}

	flagsIns, _ := flagsI.(*flags)

	return flagsIns
}
