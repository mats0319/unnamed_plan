package db

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/mats9693/unnamed_plan/public_data/config"
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

func getConn() *pg.Conn {
	return db.Conn()
}

func WithTx(task func(conn orm.DB) error) error {
	if task == nil {
		return nil // todo: return special error
	}

	conn := getConn()
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

	err = task(tx)

	return err
}

func WithNoTx(task func(conn orm.DB) error) error {
	if task == nil {
		return nil
	}

	conn := getConn()
	defer conn.Close()

	err := task(conn)

	return err
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
