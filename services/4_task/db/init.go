package db

import (
	"github.com/mats9693/unnamed_plan/services/4_task/db/dao"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"os"
)

var (
	taskDaoIns dao.TaskDao

	inited bool
)

func GetTaskDao() dao.TaskDao {
	return taskDaoIns
}

func Init() {
	if inited { // have initialized
		return
	}

	switch mdb.DB().GetDBMSName() {
	case mconst.DB_PostgreSQL:
		taskDaoIns = &dao.TaskPostgresql{}
	default:
		mlog.Logger().Error(mconst.Error_UnsupportedDB)
		os.Exit(-1)
	}

	inited = true

	mlog.Logger().Info("> Database instance init finish.")
}
