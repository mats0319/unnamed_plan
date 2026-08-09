package main

import (
	"fmt"

	mdb "github.com/mats0319/unnamed_plan/server/internal/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/db/tool_create_table/testdata"
)

func main() {
	db := mdb.InitDB(mdb.DefaultDSN, 10, 100)

	err := db.Migrator().DropTable(model.ModelList...)
	if err != nil {
		fmt.Println("drop db table failed, err: ", err)
		return
	}

	err = db.Migrator().CreateTable(model.ModelList...)
	if err != nil {
		fmt.Println("create db table failed, err: ", err)
		return
	}

	db.Create(defaultUsers)

	db.Create(testdata.TestUsers)
	db.Create(testdata.TestNotes())
	db.Create(testdata.TestFlipGameScores)
}

var defaultUsers = []*model.User{
	testdata.NewUser("mats0319", "Mario", true, false, ""),
}
