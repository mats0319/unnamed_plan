package db

import (
	"github.com/mats9693/unnamed_plan/services/1_user/db/dao"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
)

var (
	userDaoIns dao.UserDao

	inited bool
)

func GetUserDao() dao.UserDao {
	return userDaoIns
}

func Init() error {
	if inited { // have initialized
		mlog.Logger().Error("already initialized")
		return nil
	}

	switch mdb.DB().GetDBMSName() {
	case mconst.DB_PostgreSQL:
		userDaoIns = &dao.UserPostgresql{}
	default:
		mlog.Logger().Error(mconst.Error_UnsupportedDB)
		return utils.NewError(mconst.Error_UnsupportedDB)
	}

	inited = true

	mlog.Logger().Info("> Database instance init finish.")

	return nil
}
