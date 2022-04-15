package test

import (
	"context"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/init"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/test"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/mats9693/unnamed_plan/services/user/config"
	"github.com/mats9693/unnamed_plan/services/user/db"
	"github.com/mats9693/unnamed_plan/services/user/rpc"
	"go.uber.org/zap"
	"testing"
)

var (
	userServiceIns rpc_impl.IUserServer

	passed = true
)

func TestUserService(t *testing.T) {
    beforeTest()

	testLogin()
	// todo: test other func
	// for read req, only check method res is ok
	// for write req, check db record is necessary

	afterTest()
}

func testLogin() {
    // make grpc req param
    req := &rpc_impl.User_LoginReq{
        UserName: testUserName,
        Password: utils.CalcSHA256(testPassword),
    }

    // invoke method
    res, err := userServiceIns.Login(context.Background(), req)

    // check res
    if err != nil || res == nil || res.UserId != testID || res.Nickname != testNickname || res.Permission != testPermission {
        passed = false
		mlog.Logger().Error(fmt.Sprintf("> test user login failed, res: %+v\n", res), zap.Error(err))
    }
}

// beforeTest invoke 'os.exit' if anything runs unexpected
func beforeTest() {
	initialize.Init("../config_test.json", mdb.Init, config.Init, db.Init)

	// make test data
	mtest.CreateTestTable([]interface{}{
		(*model.User)(nil),
	})
	setTestUser()

	// create global service instance
	userServiceIns = rpc.GetUserServer()
}

func afterTest() {
	if passed {
		mtest.DropTestTable([]interface{}{
			(*model.User)(nil),
		})
	}
}
