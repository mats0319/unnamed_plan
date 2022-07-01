package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"net"
	"strings"
)

func FormatTarget(target ...string) ([]string, error) {
	currentIP, err := GetIP()
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(target))
	for i := range target {
		res = append(res, formatTarget(target[i], currentIP))
	}

	return res, nil
}

// formatTarget rules:
// 1. if 'target' has same ip with current service, replace target ip by '127.0.0.1'
func formatTarget(target string, ip string) (res string) {
	if strings.HasPrefix(target, ip) {
		res = "127.0.0.1" + strings.TrimPrefix(target, ip)
	} else {
		res = target
	}

	return
}

// GetIP return 192.168.2.14 ?
func GetIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", err
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}

func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return -1, err
	}

	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil // type assert is ok
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

// CompareOnStringSliceNotStrict compares two string slice, return if they have equal values ignore order
func CompareOnStringSliceNotStrict(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	valueMap := make(map[string]int, len(a)) // value - status 1: in 'a', 0: both in 'a' and 'b', -1: in 'b'
	for i := range a {
		valueMap[a[i]]++
		valueMap[b[i]]--
	}

	isEqual := true
	for _, v := range valueMap {
		if v != 0 {
			isEqual = false
			break
		}
	}

	return isEqual
}
