package mdb

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"os"
)

type dbConfig struct {
	DBMS     string `json:"DBMS"`
	Addr     string `json:"addr"`
	Database string `json:"dbName"`
	User     string `json:"user"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"` // db read and write timeout, unit: second
	ShowSQL  bool   `json:"showSQL"` // show sql before query
}

type db struct {
	config *dbConfig
	ins    DAL
}

var dbIns = &db{}

func DB() *db {
	return dbIns
}

func init() {
	dbIns.config = getDBConfig()

	switch dbIns.config.DBMS {
	case mconst.DB_PostgreSQL:
		dbIns.ins = initPostgresqlDB(dbIns.config.Addr, dbIns.config.Database, dbIns.config.User,
			dbIns.config.Password, dbIns.config.Timeout, dbIns.config.ShowSQL)
	default:
		mlog.Logger().Error("init db failed", zap.String(mconst.Error_UnsupportedDB, dbIns.config.DBMS))
		os.Exit(-1)
	}

	mlog.Logger().Info(fmt.Sprintf("> Database %s init finish.", dbIns.config.DBMS))
}

func (dbi *db) WithTx(task func(Conn) error) error {
	return dbi.ins.WithTx(task)
}

func (dbi *db) WithNoTx(task func(Conn) error) error {
	return dbi.ins.WithNoTx(task)
}

func (dbi *db) GetDBMSName() string {
	return dbi.config.DBMS
}

func getDBConfig() *dbConfig {
	byteSlice := mconfig.GetConfig(mconst.UID_DB)

	conf := &dbConfig{}
	err := json.Unmarshal(byteSlice, conf)
	if err != nil {
		mlog.Logger().Error("json unmarshal failed", zap.String("uid", mconst.UID_DB), zap.Error(err))
		return nil
	}

	return conf
}
