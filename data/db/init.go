package db

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/mats9693/unnamed_plan/config"
)

type dbConfig struct {
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"dbName"`
}

var db *pg.DB

func init() {
	conf := getDBConfig()

	db = pg.Connect(&pg.Options{
		Addr:     conf.Addr,
		User:     conf.User,
		Password: conf.Password,
		Database: conf.Database,
	})

	fmt.Println("> Database init finish.")
}

func GetDB() *pg.DB {
	return db
}

func WithTx(task func(conn orm.DB) error) error {
	if task == nil {
		return nil // todo: return special error
	}

	conn := GetDB()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			err = tx.Rollback()
		}
	}()

	return task(tx)
}

func WithNoTx(task func(conn orm.DB) error) error {
	if task == nil {
		return nil
	}

	conn := GetDB()
	defer conn.Close()

	return task(conn)
}

func getDBConfig() *dbConfig {
	byteSlice := config.GetConfig("658e06f7-71d5-4ada-b715-8c1a4489e5d2")

	conf := &dbConfig{}
	err := json.Unmarshal(byteSlice, conf)
	if err != nil {
		fmt.Printf("json unmarshal failed, uid: %s, error: %v\n", "2", err)
		return nil
	}

	return conf
}
