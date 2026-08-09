package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

var accessToken_Admin = ""
var accessToken_User = ""

var pwdSHA256 = utils.CalcSHA256("123456")

func GetAccessToken() {
	r := httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"admin","password":"%s"}`, pwdSHA256), "")
	if r.IsSuccess {
		accessToken_Admin = r.AccessToken
	}

	r = httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"user","password":"%s"}`, pwdSHA256), "")
	if r.IsSuccess {
		accessToken_User = r.AccessToken
	}
}

type TestResponse struct {
	IsSuccess   bool   `json:"is_success"`
	Err         string `json:"err"`
	AccessToken string `json:"access_token"`
	Data        struct {
		MFAToken string `json:"mfa_token"`
		TOTPKey  string `json:"totp_key"`
	} `json:"data"`
}

func httpInvoke(uri string, payload string, token string) *TestResponse {
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:10319/api"+uri, strings.NewReader(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set(utils.HTTPHeader_AccessToken, token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	r := &TestResponse{}
	err = json.Unmarshal(bodyBytes, r)
	if err != nil {
		log.Fatal(err)
	}
	if r.IsSuccess {
		// read access token
		tokenStr := res.Header.Get(utils.HTTPHeader_AccessToken)

		if len(tokenStr) > 0 {
			r.AccessToken = tokenStr
		}
	}

	return r
}

const unknownError = "unknown error"

func testCase(name string, f func() string) {
	errStr := f()
	if len(errStr) < 1 {
		log.Printf("- case: %s. Test Passed.\n", name)
	} else {
		log.Printf("- case: %s. Test Failed, error: %s\n", name, errStr)
		panic("test failed")
	}
}

// errorIs 封装错误检查，避免代码中出现大量错误文本对比，即使这个错误文本也是我写的
func errorIs(errStr string, e *utils.Error) bool {
	return errStr == e.Error()
}
