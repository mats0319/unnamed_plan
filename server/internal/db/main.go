package mdb

import (
	"log"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils/flag"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Initialize() {
	_ = InitDB(mconfig.GetInternalConfig().DBDSN, 10, 100)

	logStr := "> DB init."
	if flag.IsTestMode {
		logStr += " [test mode]"
	}
	mlog.Info(logStr)
}

var defaultDSN = "host=115.190.167.134 user=mario password=123456 dbname=test_cloud port=5432 sslmode=disable"

func Default(isTestMode bool) *gorm.DB {
	flag.IsTestMode = isTestMode

	return InitDB(defaultDSN, 10, 100)
}

func InitDB(dsn string, maxIdleConns int, maxOpenConns int) *gorm.DB {
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		SkipDefaultTransaction: true,
	}
	if flag.IsTestMode {
		gormConfig.NamingStrategy = schema.NamingStrategy{TablePrefix: "t_"}
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalln("open db failed, error: ", err)
	}

	// set conns
	{
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalln("get sql db failed, error: ", err)
		}

		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}

	dal.SetDefault(db)

	return db
}
