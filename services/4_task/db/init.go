package db

import (
	"github.com/mats9693/unnamed_plan/services/4_task/db/dao"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
)

var (
	taskDaoIns dao.TaskDao

	inited bool
)

func GetTaskDao() dao.TaskDao {
	return taskDaoIns
}

func Init() error {
	if inited { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	switch mdb.DB().GetDBMSName() {
	case mconst.DB_PostgreSQL:
		taskDaoIns = &dao.TaskPostgresql{}
	default:
		mlog.Logger().Error(mconst.Error_UnsupportedDB)
		return utils.NewError(mconst.Error_UnsupportedDB)
	}

	inited = true

	mlog.Logger().Info("> Database instance init finish.")

	return nil
}
