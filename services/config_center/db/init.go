package db

import (
    "github.com/mats9693/unnamed_plan/services/config_center/db/dao"
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "os"
)

var authenticateDaoIns dao.AuthenticateDao

func GetAuthenticateDao() dao.AuthenticateDao {
    return authenticateDaoIns
}

func init() {
    switch mdb.DB().GetDBMSName() {
    case mconst.DB_PostgreSQL:
        authenticateDaoIns = &dao.AuthenticatePostgresql{}
    default:
        mlog.Logger().Error(mconst.Error_UnsupportedDB)
        os.Exit(-1)
    }

    mlog.Logger().Info("> Database instance init finish.")
}
