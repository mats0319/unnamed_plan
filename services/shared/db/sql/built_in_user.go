package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var users = []*model.User{
	{
		UserName:   "Mario",
		Nickname:   "Mario",
		Password:   "960319", // password text
		Salt:       utils.RandomHexString(10),
		IsLocked:   false,
		Permission: 10,
		CreatedBy:  "MaTongShuai",
		Common:     model.NewCommon(),
	},
	{
		UserName:   "mario",
		Nickname:   "mario",
		Password:   "123456", // password text
		Salt:       utils.RandomHexString(10),
		IsLocked:   false,
		Permission: 10,
		CreatedBy:  "MaTongShuai",
		Common:     model.NewCommon(),
	},
	{
		UserName:   "admin",
		Nickname:   "admin",
		Password:   "admin", // password text
		Salt:       utils.RandomHexString(10),
		IsLocked:   false,
		Permission: 6,
		CreatedBy:  "MaTongShuai",
		Common:     model.NewCommon(),
	},
}

func setDefaultUser(db *pg.DB) {
	for i := range users {
		users[i].Password = utils.CalcSHA256(users[i].Password)
		users[i].Password = utils.CalcSHA256(users[i].Password, users[i].Salt)

		_, err := db.Model(users[i]).Insert()
		if err != nil {
			log.Printf("set default users failed, index: %d, error: %v\n", i, err)
		}
	}
}
