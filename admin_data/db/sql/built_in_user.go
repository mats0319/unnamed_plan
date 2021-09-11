package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/utils/uuid"
	"log"
)

var users = []*model.User{
	{
		UserID:     uuid.New(),
		Name:       "Mario",
		Nickname:   "Mario",
		Password:   "960319",
		IsLocked:   false,
		Permission: 10,
		Common:     model.NewCommon(),
	},
	{
		UserID:     uuid.New(),
		Name:       "admin",
		Nickname:   "admin",
		Password:   "admin",
		IsLocked:   false,
		Permission: 6,
		Common:     model.NewCommon(),
	},
}

func insertUsers(db *pg.DB) {
	for i := range users {
		_, err := db.Model(users[i]).Insert()
		if err != nil {
			log.Fatalf("insert users failed, index: %d, error: %v\n", i, err)
		}
	}
}
