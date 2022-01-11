package db

import (
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/user/db/dao"
    "os"
)

var userDaoIns dao.UserDao

func GetUserDao() dao.UserDao {
    return userDaoIns
}

func init() {
    switch mdb.DB().GetDBMSName() {
    case mconst.DB_PostgreSQL:
        userDaoIns = &dao.UserPostgresql{}
    default:
        mlog.Logger().Error(mconst.Error_UnsupportedDB)
        os.Exit(-1)
    }

    mlog.Logger().Info("> Database instance init finish.")
}
