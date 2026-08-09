package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand/v2"
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

const charactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示字符库中的全部字符

// GenerateRandomBytes generate random readable Bytes
func GenerateRandomBytes[T string | []byte](length int) T {
	b := make([]byte, length)

	randomNum, remainBits := rand.Int64(), 64
	for i := 0; i < len(b); {
		if remainBits < useBits {
			randomNum, remainBits = rand.Int64(), 64
		}

		index := int(randomNum & (1<<useBits - 1)) // 0b0011 1111
		if index < len(charactersLibrary) {
			randomNum >>= useBits
			remainBits -= useBits

			b[i] = charactersLibrary[index]
			i++
		} else {
			randomNum >>= 1
			remainBits -= 1
		}
	}

	return T(b)
}
