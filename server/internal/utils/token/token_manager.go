package token

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type TokenManager struct {
	HMACKey string
}

var tm = &TokenManager{}

func InitTokenManager(hmacKey string) {
	tm.HMACKey = hmacKey
}

func SerializeToken(token *Token) string {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		e := utils.ErrServerInternalError().WithCause(err)
		mlog.Error(e.String())
		return ""
	}

	tokenHex := hex.EncodeToString(tokenBytes)
	tokenHash := utils.HMACSHA256(tokenHex, tm.HMACKey)

	return fmt.Sprintf("%s.%s", tokenHex, tokenHash)
}

// DeserializeToken no log
func DeserializeToken(token string, typ TokenType) (t *Token, e *utils.Error) {
	// check token structure and hash
	tokenSplit := strings.Split(token, ".")
	if len(tokenSplit) != 2 {
		e = utils.ErrInvalidToken().WithParam("token", token)
		return
	}

	if utils.HMACSHA256(tokenSplit[0], tm.HMACKey) != tokenSplit[1] { // check token hash
		e = utils.ErrWrongTokenHash().WithParam("payload", tokenSplit[0])
		return
	}

	// decode & deserialize
	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e = utils.ErrDecodeToken().WithCause(err).WithParam("payload", tokenSplit[0])
		return
	}

	t = &Token{}
	err = json.Unmarshal(tokenBytes, t)
	if err != nil {
		e = utils.ErrDeserializeToken().WithCause(err)
		return
	}

	// check token 'type' and 'expire time'
	if t.Type != typ {
		e = utils.ErrInvalidTokenType().WithParam("want", typ).WithParam("get", t.Type)
		return
	}

	if t.ExpireTime < time.Now().UnixMilli() {
		e = utils.ErrTokenExpired().WithParam("expire time", t.ExpireTime)
		return
	}

	return
}
