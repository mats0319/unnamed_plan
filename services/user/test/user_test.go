package test

import (
	"context"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/shared/db"
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
	beforeTest_Postgresql(t)

	testLogin()
	testList()
	testCreate()
	testLock()
	testUnlock()
	testModifyInfo()
	testModifyPermission()

	afterTest_Postgresql(t)
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

func testList() {
	// make grpc req param
	req := &rpc_impl.User_ListReq{
		OperatorId: testID,
		Page: &rpc_impl.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
	}

	// invoke method
	res, err := userServiceIns.List(context.Background(), req)

	// check res
	if err != nil || res == nil || res.Total != 0 {
		passed = false
		mlog.Logger().Error(fmt.Sprintf("> test user list failed, res: %+v\n", res), zap.Error(err))
	}
}

func testCreate() {
	// make grpc req param
	req := &rpc_impl.User_CreateReq{
		OperatorId: testID,
		UserName:   testCreateUserName,
		Password:   utils.CalcSHA256(testPassword),
		Permission: testCreatePermission,
	}

	// invoke method
	_, err := userServiceIns.Create(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName()

	if err != nil || err2 != nil || newUser.Permission != testCreatePermission {
		passed = false
		mlog.Logger().Error("> test user create failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}

	// for next steps: get 'user id' of create user
	testCreateUserID = newUser.ID
}

func testLock() {
	// make grpc req param
	req := &rpc_impl.User_LockReq{
		OperatorId: testID,
		UserId:     testCreateUserID,
	}

	// invoke method
	_, err := userServiceIns.Lock(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName()

	if err != nil || err2 != nil || !newUser.IsLocked {
		passed = false
		mlog.Logger().Error("> test user lock failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func testUnlock() {
	// make grpc req param
	req := &rpc_impl.User_UnlockReq{
		OperatorId: testID,
		UserId:     testCreateUserID,
	}

	// invoke method
	_, err := userServiceIns.Unlock(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName()

	if err != nil || err2 != nil || newUser.IsLocked {
		passed = false
		mlog.Logger().Error("> test user unlock failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func testModifyInfo() {
	// make grpc req param
	req := &rpc_impl.User_ModifyInfoReq{
		OperatorId: testCreateUserID,
		UserId:     testCreateUserID,
		CurrPwd:    utils.CalcSHA256(testPassword),
		Nickname:   testModifyInfoUserName,
		Password:   "",
	}

	// invoke method
	_, err := userServiceIns.ModifyInfo(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName()

	if err != nil || err2 != nil || newUser.Nickname != testModifyInfoUserName {
		passed = false
		mlog.Logger().Error("> test user modify info failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func testModifyPermission() {
	// make grpc req param
	req := &rpc_impl.User_ModifyPermissionReq{
		OperatorId: testID,
		UserId:     testCreateUserID,
		Permission: testModifyPermissionP,
	}

	// invoke method
	_, err := userServiceIns.ModifyPermission(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName()

	if err != nil || err2 != nil || newUser.Permission != testModifyPermissionP {
		passed = false
		mlog.Logger().Error("> test user modify permission failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func getUserByUserName() (*model.User, error) {
	newUser := &model.User{}
	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(newUser).Where(model.User_UserName+" = ?", testCreateUserName).First()
	})

	return newUser, err
}

func beforeTest_Postgresql(t *testing.T) {
	initialize.Init("../config_test.json", mdb.Init, config.Init, db.Init)

	// create test table
	err := mtest.CreateTestTable_Postgresql([]interface{}{
		(*model.User)(nil),
	})
	if err != nil {
		t.Fail()
	}

	// pre-set test data
	setTestUser()

	// create global service instance
	userServiceIns = rpc.GetUserServer()
}

func afterTest_Postgresql(t *testing.T) {
	if passed {
		mtest.DropTestTable_Postgresql([]interface{}{
			(*model.User)(nil),
		})
	} else {
		t.Fail()
	}
}
