package rpc

import (
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "github.com/mats9693/unnamed_plan/services/user/db"
    "go.uber.org/zap"
    "os"
)

type testData struct {
    testUsers []*model.User

    testUserName string
    testNickname   string
    testPassword   string
    testPermission uint8
    testID         string

    testCreateUserName   string
    testCreatePermission uint8
    testCreateUserID     string // var

    testModifyInfoUserName string

    testModifyPermissionP uint8
}

func (td *testData) loadTestData() {
    td.testUserName = "test user"
    td.testNickname = "test nickname"
    td.testPassword = "test password"
    td.testPermission = 8
    td.testID = "test id"

    td.testCreateUserName = "test create user"
    td.testCreatePermission = 1

    td.testModifyInfoUserName = "test create user - modify"
    td.testModifyPermissionP = 2

    td.testUsers = []*model.User{
        {
            UserName:   td.testUserName,
            Nickname:   td.testNickname,
            Password:   td.testPassword, // password text
            Salt:       utils.RandomHexString(10),
            IsLocked:   false,
            Permission: td.testPermission,
            CreatedBy:  "MaTongShuai",
            Common: model.Common{
                ID: td.testID,
            },
        },
    }

    setTestUser(td.testUsers)
}

func setTestUser(users []*model.User) {
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
