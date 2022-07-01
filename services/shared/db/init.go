package mdb

import (
	"encoding/json"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/config"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
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

	dal mdal.DAL
}

var dbIns = &db{}

func DB() *db {
	return dbIns
}

func Init() error {
	if dbIns.config != nil { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	dbIns.config = getDBConfig()

	if dbIns.config == nil {
		mlog.Logger().Error("config not ready when init db")
		return utils.NewError("config not ready when init db")
	}

	switch dbIns.config.DBMS {
	case mconst.DB_PostgreSQL:
		dbIns.dal = mdal.InitPostgresqlDB(dbIns.config.Addr, dbIns.config.Database, dbIns.config.User,
			dbIns.config.Password, dbIns.config.Timeout, dbIns.config.ShowSQL)
	default:
		mlog.Logger().Error("init db failed", zap.String(mconst.Error_UnsupportedDB, dbIns.config.DBMS))
		return utils.NewError(mconst.Error_UnsupportedDB + dbIns.config.DBMS)
	}

	mlog.Logger().Info(fmt.Sprintf("> DBMS %s init finish.", dbIns.config.DBMS))

	return nil
}

func (dbi *db) WithTx(task func(mdal.Conn) error) error {
	return dbi.dal.WithTx(task)
}

func (dbi *db) WithNoTx(task func(mdal.Conn) error) error {
	return dbi.dal.WithNoTx(task)
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
