package mfa

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"slices"
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func VerifyTOTPCode[T string | []byte](code string, keyBase32 T) *utils.Error {
	if len(code) != 6 {
		e := utils.ErrInvalidTOTPCode().WithParam("code", code)
		mlog.Error(e.String())
		return e
	}

	key := make([]byte, 10)
	n, err := base32.StdEncoding.Decode(key, []byte(keyBase32))
	if err != nil {
		e := utils.ErrInvalidTOTPKey().WithCause(err).WithParam("base32 str", keyBase32)
		mlog.Error(e.String())
		return e
	}
	key = key[:n]

	timestep := time.Now().Unix() / 30

	// allow last/current/next totp code
	validCodes := []string{
		calcTOTPCode(key, iTob(timestep-1)),
		calcTOTPCode(key, iTob(timestep)),
		calcTOTPCode(key, iTob(timestep+1)),
	}

	if !slices.Contains(validCodes, code) {
		e := utils.ErrWrongTOTPCode().WithParam("code", code)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func calcTOTPCode(key []byte, content []byte) string {
	hasher := hmac.New(sha1.New, key)
	hasher.Write(content)
	hmacHash := hasher.Sum(nil)

	offset := int(hmacHash[len(hmacHash)-1] & 0x0f)
	// 算法要求屏蔽最高有效位
	longPassword := int(hmacHash[offset]&0x7f)<<24 |
		int(hmacHash[offset+1])<<16 |
		int(hmacHash[offset+2])<<8 |
		int(hmacHash[offset+3])

	password := longPassword % int(math.Pow10(6))

	return fmt.Sprintf("%06d", password)
}

func iTob(integer int64) []byte {
	byteSlice := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteSlice[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteSlice
}
