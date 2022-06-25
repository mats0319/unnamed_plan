package db

import (
	"github.com/mats9693/unnamed_plan/services/2_cloud_file/db/dao"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
)

var (
	cloudFileDaoIns dao.CloudFileDao

	inited bool
)

func GetCloudFileDao() dao.CloudFileDao {
	return cloudFileDaoIns
}

func Init() error {
	if inited { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	switch mdb.DB().GetDBMSName() {
	case mconst.DB_PostgreSQL:
		cloudFileDaoIns = &dao.CloudFilePostgresql{}
	default:
		mlog.Logger().Error(mconst.Error_UnsupportedDB)
		return utils.NewError(mconst.Error_UnsupportedDB)
	}

	inited = true

	mlog.Logger().Info("> Database instance init finish.")

	return nil
}
