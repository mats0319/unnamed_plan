package mfa

import (
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/token"
)

func GenerateMFAToken(userName string, expireMinute int) string {
	return token.SerializeToken(&token.Token{
		UserName:   userName,
		Type:       token.TokenType_MFAToken,
		ExpireTime: time.Now().Add(time.Duration(expireMinute) * time.Minute).UnixMilli(),
	})
}

func VerifyMFAToken(tokenStr string) (t *token.Token, e *utils.Error) {
	t, e = token.DeserializeToken(tokenStr, token.TokenType_MFAToken)
	if e != nil {
		mlog.Error(e.String())
		return
	}
	
	return
}
