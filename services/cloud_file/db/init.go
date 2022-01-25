package db

import (
    "github.com/mats9693/unnamed_plan/services/cloud_file/db/dao"
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "os"
)

var cloudFileDaoIns dao.CloudFileDao

func GetCloudFileDao() dao.CloudFileDao {
    return cloudFileDaoIns
}

func init() {
    switch mdb.DB().GetDBMSName() {
    case mconst.DB_PostgreSQL:
        cloudFileDaoIns = &dao.CloudFilePostgresql{}
    default:
        mlog.Logger().Error(mconst.Error_UnsupportedDB)
        os.Exit(-1)
    }

    mlog.Logger().Info("> Database instance init finish.")
}