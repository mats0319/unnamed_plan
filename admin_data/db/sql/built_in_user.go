package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/utils/uuid"
	"log"
	"math/rand"
)

var users = []*model.User{
	{
		UserID:     uuid.New(),
		UserName:   "Mario",
		Nickname:   "Mario",
		Password:   "960319", // password text
		Salt:       randomString(10),
		IsLocked:   false,
		Permission: 10,
		Common:     model.NewCommon(),
	},
	{
		UserID:     uuid.New(),
		UserName:   "admin",
		Nickname:   "admin",
		Password:   "admin", // password text
		Salt:       randomString(10),
		IsLocked:   false,
		Permission: 6,
		Common:     model.NewCommon(),
	},
}

func insertUsers(db *pg.DB) {
	for i := range users {
		users[i].Password = calcPassword(users[i].Password, users[i].Salt)

		_, err := db.Model(users[i]).Insert()
		if err != nil {
			log.Fatalf("insert users failed, index: %d, error: %v\n", i, err)
		}
	}
}

func calcPassword(text string, salt string) string {
	hash := sha256.New()
	hash.Reset()
	hash.Write([]byte(text + salt))
	bytes := hash.Sum(nil)

	return hex.EncodeToString(bytes)
}

const str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = str[rand.Intn(len(str))]
	}

	return string(bytes)
}
