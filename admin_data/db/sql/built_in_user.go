package main

import (
	"log"

	"github.com/go-pg/pg/v10"

	"github.com/mats9693/unnamed_plan/admin_data/db/model"
    "github.com/mats9693/unnamed_plan/admin_data/utils"
	"github.com/mats9693/utils/uuid"
)

var users = []*model.User{
	{
		UserID:     uuid.New(),
		UserName:   "Mario",
		Nickname:   "Mario",
		Password:   "960319", // password text
		Salt:       utils.RandomString(10),
		IsLocked:   false,
		Permission: 10,
		CreatedBy:  "MaTongShuai",
		Common:     model.NewCommon(),
	},
	{
		UserID:     uuid.New(),
		UserName:   "admin",
		Nickname:   "admin",
		Password:   "admin", // password text
		Salt:       utils.RandomString(10),
		IsLocked:   false,
		Permission: 6,
		CreatedBy:  "MaTongShuai",
		Common:     model.NewCommon(),
	},
}

func insertUsers(db *pg.DB) {
	for i := range users {
		users[i].Password = utils.CalcSHA256(users[i].Password)
		users[i].Password = utils.CalcSHA256(users[i].Password, users[i].Salt)

		_, err := db.Model(users[i]).Insert()
		if err != nil {
			log.Printf("insert users failed, index: %d, error: %v\n", i, err)
		}
	}
}
