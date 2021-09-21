package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/unnamed_plan/admin_data/kits"
	"github.com/mats9693/utils/uuid"
	"log"
)

var users = []*model.User{
	{
		UserID:     uuid.New(),
		UserName:   "Mario",
		Nickname:   "Mario",
		Password:   "960319", // password text
		Salt:       kits.RandomString(10),
		IsLocked:   false,
		Permission: 10,
		Common:     model.NewCommon(),
	},
	{
		UserID:     uuid.New(),
		UserName:   "admin",
		Nickname:   "admin",
		Password:   "admin", // password text
		Salt:       kits.RandomString(10),
		IsLocked:   false,
		Permission: 6,
		Common:     model.NewCommon(),
	},
}

func insertUsers(db *pg.DB) {
	for i := range users {
		users[i].Password = kits.CalcPassword(users[i].Password, users[i].Salt)

		_, err := db.Model(users[i]).Insert()
		if err != nil {
			log.Fatalf("insert users failed, index: %d, error: %v\n", i, err)
		}
	}
}
