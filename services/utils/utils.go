package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"strings"
)

// CalcSHA256 calc sha256('text'[+'extension'])
func CalcSHA256(text string, append ...string) string {
	hash := sha256.New()
	hash.Reset()
	hash.Write([]byte(text + strings.Join(append, "")))
	bytes := hash.Sum(nil)

	return hex.EncodeToString(bytes)
}

func StringToBool(str string) (res bool, err error) {
	switch str {
	case "true":
		res = true
	case "false":
		res = false
	default:
		err = errors.New("unknown str:" + str)
	}

	return
}

func ErrorsToString(errs ...error) string {
	res := ""
	for i := range errs {
		if errs[i] != nil {
			res += errs[i].Error()
		}
	}

	return res
}
