package db

import (
	"github.com/mats9693/unnamed_plan/services/config_center/db/dao"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
)

var (
	serviceConfigDaoIns dao.ServiceConfigDao
	configItemDaoIns    dao.ConfigItemDao

	inited bool
)

func GetServiceConfigDao() dao.ServiceConfigDao {
	return serviceConfigDaoIns
}

func GetConfigItemDao() dao.ConfigItemDao {
	return configItemDaoIns
}

func Init() error {
	if inited { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	switch mdb.DB().GetDBMSName() {
	case mconst.DB_PostgreSQL:
		serviceConfigDaoIns = &dao.ServiceConfigPostgresql{}
		configItemDaoIns = &dao.ConfigItemPostgresql{}
	default:
		mlog.Logger().Error(mconst.Error_UnsupportedDB)
		return utils.NewError(mconst.Error_UnsupportedDB)
	}

	inited = true

	mlog.Logger().Info("> Database instance init finish.")

	return nil
}
