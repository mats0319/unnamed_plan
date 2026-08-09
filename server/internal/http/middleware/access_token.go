package middleware

import (
	"time"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/token"
)

func GenerateAPIAccessToken(userName string, expireHour int) string {
	return token.SerializeToken(&token.Token{
		UserName:   userName,
		Type:       token.TokenType_APIAccessToken,
		ExpireTime: time.Now().Add(time.Duration(expireHour) * time.Hour).UnixMilli(),
	})
}

// OptionalVerifyAPIAccessToken 用于访客和用户均可使用的接口，访问token验证失败视为访客、验证成功可以拿到`ctx.userName`。
func OptionalVerifyAPIAccessToken(ctx *mhttp.Context) *utils.Error {
	_ = verifyAPIAccessToken(ctx)
	return nil
}

func VerifyAPIAccessToken(ctx *mhttp.Context) *utils.Error {
	e := verifyAPIAccessToken(ctx)
	if e != nil {
		mlog.Error(e.String())
		return utils.ErrInvalidAccessToken()
	}

	return nil
}

// verifyAccessToken 函数内不打印错误，因为部分场景允许验证错误（例如上方的可选验证函数）
func verifyAPIAccessToken(ctx *mhttp.Context) *utils.Error {
	t, e := token.DeserializeToken(ctx.AccessToken, token.TokenType_APIAccessToken)
	if e != nil {
		return e
	}

	ctx.UserName = t.UserName

	return nil
}
