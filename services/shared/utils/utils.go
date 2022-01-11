package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
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

func NewError(data string) error {
	return errors.New(data)
}

// FormatDirSuffix make sure path is directory(end with '/')
func FormatDirSuffix(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path
}

func RandomHexString(length int) string {
	str := ""
	for len(str) < length {
		str += fmt.Sprintf("%x", rand.Uint64())
	}

	return str[:length]
}

func Contains(slice []string, value string) bool {
	isValid := false
	for i := range slice {
		if value == slice[i] {
			isValid = true
			break
		}
	}

	return isValid
}

func GetIndex(slice []string, value string) int {
	index := -1
	for i := range slice {
		if value == slice[i] {
			index = i
			break
		}
	}

	return index
}
