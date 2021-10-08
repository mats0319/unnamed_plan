package kits

import (
	"fmt"
	"testing"
)

func TestCalcPassword(t *testing.T) {
	pwdHash := CalcSHA256("123456")

	fmt.Println("> --- Test calc password ---")
	fmt.Println("> password hash :", pwdHash)
	fmt.Println("> final password:", CalcSHA256(pwdHash, "zdoZPfZxsT"))
}
