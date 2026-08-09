package main

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/test/api"
)

// start server with flag '-t' to use test db connection
func main() {
	log.Println("> Test Start ...")

	createTable()

	api.GetAccessToken()

	// todo: 编写更多异常用例
	testApi("Register", api.Register)
	testApi("Login", api.Login) // include enable/disable MFA
	testApi("List User", api.ListUser)
	testApi("Modify User", api.ModifyUser)
	testApi("New TOTP Key", api.NewTOTPKey)
	testApi("Set MFA Status", api.SetMFAStatus)

	testApi("Create Note", api.CreateNote)
	testApi("List Note", api.ListNote)
	testApi("Modify Note", api.ModifyNote)
	testApi("Delete Note", api.DeleteNote)

	testApi("List Game Score", api.ListGameScore)
	testApi("Upload Game Score", api.UploadGameScore)

	//dropTable() // do not del testdata during dev

	log.Println("> All Test Passed! ^_^")
}

func testApi(name string, f func()) {
	log.Printf("> %s.\n", name)
	defer log.Println()

	f()
}
