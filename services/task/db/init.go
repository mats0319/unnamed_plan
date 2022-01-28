package db

import (
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/task/db/dao"
    "os"
)

var taskDaoIns dao.TaskDao

func GetTaskDao() dao.TaskDao {
    return taskDaoIns
}

func init() {
    switch mdb.DB().GetDBMSName() {
    case mconst.DB_PostgreSQL:
        taskDaoIns = &dao.TaskPostgresql{}
    default:
        mlog.Logger().Error(mconst.Error_UnsupportedDB)
        os.Exit(-1)
    }

    mlog.Logger().Info("> Database instance init finish.")
}
