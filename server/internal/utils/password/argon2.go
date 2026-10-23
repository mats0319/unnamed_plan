package password

import (
	"crypto/subtle"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

// 数据库密码存储不应使用明文，前端传参数`pwdSHA256`、服务端保存的则是`pwdArgon2`：
//
// - `pwdSHA256`: `hex(sha256('password text'))`，服务端接收该结构的参数
// - `key`: `argon2.derive('pwdSHA256', 'salt')`，算法派生的密钥
//     - argon2算法可以根据算法参数、密码和盐派生密钥（相同环境下，重复派生的结果相同）
//     - `salt`: 程序随机生成
// - `pwdArgon2`/`pwdDB`: `encode(key)`，数据库保存的密码
//     - `encode`：按照一定格式，将密码和相关信息编码成字符串
//     - 编码结构: `argon2id.v=19,m=65536,t=3,p=1.[saltHex].[keyHex]`
//
// 前端执行hash的目的是即使恶意攻击窃取了http请求参数，也无法恢复出密码明文，进而攻击者无法使用我们的UI；
// 后端执行hash的目的是即使恶意攻击窃取了数据库数据，也无法恢复出密码hash，进而无法调用接口。

func GeneratePassword(pwdSHA256 string) (pwdArgon2 string) {
	params := defaultAlgorithmParams()

	key := params.deriveKey(pwdSHA256)
	pwdArgon2 = params.encode(key)

	return
}

func VerifyPassword(pwdSHA256 string, pwdArgon2 string) *utils.Error {
	params, oldKey, e := (&AlgorithmParams{}).decode(pwdArgon2)
	if e != nil {
		return e
	}

	newKey := params.deriveKey(pwdSHA256)

	if subtle.ConstantTimeCompare(oldKey, newKey) != 1 { // 使用恒定时间比较防止时序攻击
		// 这里不打印错误，因为部分应用场景要求密码验证不能通过（例如修改密码时，新、旧密码不能一样）。
		// 具体的，假设这里打印错误，那么在修改密码时，即使一切正确执行，控制台也会报错。
		// 虽然这样可能导致前面decode出错时，错误被打印2次，但我认为这是可以接受的
		return utils.ErrWrongPassword()
	}

	return nil
}
