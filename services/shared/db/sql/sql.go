package main

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"log"
	"math/rand"
	"reflect"
	"time"
)

var dbConf = &pg.Options{
	Addr:         "106.75.174.54:5432",
	User:         "mario_test",
	Password:     "123456",
	Database:     "unnamed_plan_test",
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	db := pg.Connect(dbConf)
	defer func() {
		_ = db.Close()
	}()

	models := []interface{}{
		(*model.User)(nil),
		(*model.CloudFile)(nil),
		(*model.Note)(nil),
		(*model.Task)(nil),
	}

	for i, m := range models {
		err := db.Model(m).CreateTable(&orm.CreateTableOptions{Temp: false, IfNotExists: true})
		if err != nil {
			log.Fatalf("create table index: %d, name: %s failed, error: %v\n", i, reflect.TypeOf(m), err)
		}
	}

	// built-in data
	{
		setDefaultUser(db)
	}

	log.Println("db init finish.")
}
