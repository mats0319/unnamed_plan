package mfa

import (
	"encoding/base32"
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func GenerateTOTPKey(userName string, expireMinute int, encryptKey string) (key string, e *utils.Error) {
	keyBytes := utils.GenerateRandomBytes[[]byte](10)
	keyBase32 := base32.StdEncoding.EncodeToString(keyBytes)
	keyEncrypted, err := utils.Encrypt(keyBase32, encryptKey)
	if err != nil {
		e = utils.ErrEncrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	expireTime := time.Now().Add(time.Duration(expireMinute) * time.Minute).UnixMilli()

	newTOTPKey(userName, keyEncrypted, expireTime)

	key = keyBase32

	return
}

func VerifyTOTPKey(userName string, code string, encryptKey string) (key string, e *utils.Error) {
	return maintainTOTPKey(userName, code, encryptKey)
}
