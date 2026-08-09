package main

import (
	"log"

	mdb "github.com/mats0319/unnamed_plan/server/internal/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/test/testdata"
	"gorm.io/gorm"
)

var db *gorm.DB

func createTable() {
	dropTable()

	err := db.Migrator().CreateTable(model.ModelList...)
	if err != nil {
		log.Fatalln("create db table failed, error: ", err)
	}

	db.Create(testdata.PresetUser())
}

func dropTable() {
	if db == nil {
		db = mdb.InitTestDB()
	}

	err := db.Migrator().DropTable(model.ModelList...)
	if err != nil {
		log.Fatalln("drop db table failed, error: ", err)
	}
}
