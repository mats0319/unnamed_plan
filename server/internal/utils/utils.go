package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

// HMACSHA256 calc hmac-sha256('key', 'content'), return hex(hash)
func HMACSHA256[T string | []byte](content string, key T) string {
	hasher := hmac.New(sha256.New, []byte(key)) // k default nil is ok
	hasher.Write([]byte(content))
	bytes := hasher.Sum(nil)

	return hex.EncodeToString(bytes)
}

func CalcSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	bytes := hasher.Sum(nil)

	return hex.EncodeToString(bytes)
}

func UUIDv5[T string | []byte](data T) string {
	return strings.ToUpper(uuid.NewSHA1(uuid.NameSpaceDNS, []byte(data)).String())
}

func Encrypt[T string | []byte](message T, key T) (cipherHex string, e *Error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		e = ErrEncrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	aesgcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		e = ErrEncrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	cipherBytes := aesgcm.Seal(nil, nil, []byte(message), nil)
	cipherHex = hex.EncodeToString(cipherBytes)

	return
}

func Decrypt[T string | []byte](cipherHex string, key T) (message []byte, e *Error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		e = ErrDecrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	aesgcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		e = ErrDecrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	cipherBytes, err := hex.DecodeString(cipherHex)
	if err != nil {
		e = ErrDecrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	message, err = aesgcm.Open(nil, nil, cipherBytes, nil)
	if err != nil {
		e = ErrDecrypt().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}

func GenerateRandomBytes(length int) []byte {
	bytesBuilder := bytes.NewBuffer(nil)

	l := length
	for l > 0 {
		_, _ = bytesBuilder.WriteString(rand.Text()) // err always nil
		l -= 26
	}

	return []byte(bytesBuilder.String())[:length]
}
