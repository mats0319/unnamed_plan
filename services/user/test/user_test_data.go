package test

import (
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/mats9693/unnamed_plan/services/user/db"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const (
	testUserName = "test user"
	testNickname = "test nickname"
	testPassword = "test password"
	testPermission = 1
	testID = "test id"
)

var users = []*model.User{
	{
		UserName:   testUserName,
		Nickname:   testNickname,
		Password:   testPassword, // password text
		Salt:       utils.RandomHexString(10),
		IsLocked:   false,
		Permission: testPermission,
		CreatedBy:  "MaTongShuai",
		Common: model.Common{
			ID: testID,
		},
	},
}

func setTestUser() {
	for i := range users {
		users[i].Password = utils.CalcSHA256(users[i].Password)
		users[i].Password = utils.CalcSHA256(users[i].Password, users[i].Salt)

		err := db.GetUserDao().Insert(users[i])
		if err != nil {
			mlog.Logger().Error("set default users failed", zap.Int("index", i), zap.Error(err))
			os.Exit(-1)
		}
	}
}
