package password

import (
	"encoding/hex"
	"fmt"
	"strings"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"golang.org/x/crypto/argon2"
)

type AlgorithmParams struct {
	CalcTimes uint32 // 迭代次数
	Memory    uint32 // 使用内存
	Threads   uint8  // 使用线程数
	KeyLength uint32
	Salt      []byte
}

func defaultAlgorithmParams() *AlgorithmParams {
	return &AlgorithmParams{
		CalcTimes: 3,
		Memory:    64 * 1024, // 64 MB
		Threads:   1,
		KeyLength: 32,
		Salt:      utils.GenerateRandomBytes[[]byte](32),
	}
}

func (ap *AlgorithmParams) deriveKey(pwdSHA256 string) []byte {
	if ap == nil {
		return make([]byte, 32)
	}

	return argon2.IDKey([]byte(pwdSHA256), ap.Salt, ap.CalcTimes, ap.Memory, ap.Threads, ap.KeyLength)
}

func (ap *AlgorithmParams) encode(key []byte) string {
	if ap == nil {
		return ""
	}

	saltHex := hex.EncodeToString(ap.Salt)
	keyHex := hex.EncodeToString(key)

	return fmt.Sprintf("argon2id.v=%d,m=%d,t=%d,c=%d.%s.%s",
		argon2.Version, ap.Memory, ap.CalcTimes, ap.Threads, saltHex, keyHex)
}

func (ap *AlgorithmParams) decode(pwdArgon2 string) (params *AlgorithmParams, oldKey []byte, e *utils.Error) {
	pwdSplit := strings.Split(pwdArgon2, ".")
	if len(pwdSplit) != 4 || pwdSplit[0] != "argon2id" {
		e = utils.ErrInvalidPassword().WithParam("encoded pwd", pwdArgon2)
		mlog.Error(e.String())
		return
	}

	params = &AlgorithmParams{}

	var version int
	_, err := fmt.Sscanf(pwdSplit[1], "v=%d,m=%d,t=%d,c=%d", &version, &params.Memory, &params.CalcTimes, &params.Threads)
	if err != nil || version != argon2.Version {
		e = utils.ErrInvalidPassword().WithCause(err).WithParam("params", params).WithParam("version", version)
		mlog.Error(e.String())
		return
	}

	params.Salt, err = hex.DecodeString(pwdSplit[2])
	if err != nil {
		e = utils.ErrInvalidPwdSalt().WithCause(err).WithParam("salt", pwdSplit[2])
		mlog.Error(e.String())
		return
	}

	oldKey, err = hex.DecodeString(pwdSplit[3])
	if err != nil {
		e = utils.ErrInvalidPwdKey().WithCause(err).WithParam("key", pwdSplit[3])
		mlog.Error(e.String())
		return
	}
	params.KeyLength = uint32(len(oldKey))

	return
}
