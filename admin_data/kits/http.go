package kits

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"time"

	"github.com/mats9693/unnamed_plan/admin_data/db/model"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func SortUsersByUserID(users []*model.User, order []string) ([]*model.User, error) {
	if len(users) != len(order) {
		return nil, errors.New(fmt.Sprintf("unmatched users amount, users %d, orders %d", len(users), len(order)))
	}

	length := len(users)
	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			if order[j] == users[i].UserID {
				users[i], users[j] = users[j], users[i]
				break
			}
		}
	}

	unmatchedIndex := -1
	for i := 0; i < length; i++ {
		if users[i].UserID != order[i] {
			unmatchedIndex = i
			break
		}
	}

	if unmatchedIndex >= 0 {
		return nil, errors.New(fmt.Sprintf("unmatched user id: %s", users[unmatchedIndex].UserID))
	}

	return users, nil
}

const str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = str[rand.Intn(len(str))]
	}

	return string(bytes)
}

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

// AppendDirSuffix make sure path is directory(end with '/')
func AppendDirSuffix(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path
}
