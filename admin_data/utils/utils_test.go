package utils

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/admin_data/http/structure_defination"
	"testing"
)

func TestCalcPassword(t *testing.T) {
	pwdHash := CalcSHA256("123456")

	fmt.Println("> --- Test calc password ---")
	fmt.Println("> password hash :", pwdHash)
	fmt.Println("> final password:", CalcSHA256(pwdHash, "zdoZPfZxsT"))
}

// TestJSONMarshalLevel 结论：使用匿名结构体嵌套，json marshal不会引入新的一层，匿名的指针结构体和值结构体都不会
func TestJSONMarshalLevel(t *testing.T) {
	type Inner struct {
		Name string `json:"name"`
	}

	wrapper := &struct {
		InnerIns Inner
		*Inner
		Password string `json:"pwd"`
	}{
		InnerIns: Inner{Name: "wrapper"},
		Inner: &Inner{Name: "direct"},
		Password: "123",
	}

	wrapperByte, _ := json.Marshal(wrapper)
	fmt.Println("> wrapper: " + string(wrapperByte))

	// > wrapper: {"InnerIns":{"name":"wrapper"},"name":"direct","pwd":"123"}
}

// TestJSONMarshal 确认包装功能
func TestJSONMarshal(t *testing.T) {
	data := structure.MakeLoginRes("id", "nickname", 1)
	dataByte, _ := json.Marshal(data)
	fmt.Println(string(dataByte))

	// {"userID":"id","nickname":"nickname","permission":1}
}
