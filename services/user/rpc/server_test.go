package rpc

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
	"go.uber.org/zap"
	"testing"
)

type userServiceTest struct {
	service rpc_impl.IUserServer
	passed  bool // if all test passed

	testData *testData
}

// TestUserService use 'postgresql' db, if you want test in another db, you need to write matched test code
func TestUserService(t *testing.T) {
	serviceIns := &userServiceTest{testData: &testData{}}

	serviceIns.beforeTest(t)

	serviceIns.testLogin()
	serviceIns.testList()
	serviceIns.testCreate()
	serviceIns.testLock()
	serviceIns.testUnlock()
	serviceIns.testModifyInfo()
	serviceIns.testModifyPermission()

	serviceIns.afterTest(t)
}

func (s *userServiceTest) testLogin() {
	// make grpc req param
	req := &rpc_impl.User_LoginReq{
		UserName: s.testData.testUserName,
		Password: utils.CalcSHA256(s.testData.testPassword),
	}

	// invoke method
	res, err := s.service.Login(context.Background(), req)

	// check res
	if err != nil || res == nil || res.UserId != s.testData.testID || res.Nickname != s.testData.testNickname ||
		res.Permission != uint32(s.testData.testPermission) {
		s.passed = false
		mlog.Logger().Error(fmt.Sprintf("> test user login failed, res: %+v\n", res), zap.Error(err))
	}
}

func (s *userServiceTest) testList() {
	// make grpc req param
	req := &rpc_impl.User_ListReq{
		OperatorId: s.testData.testID,
		Page: &rpc_impl.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
	}

	// invoke method
	res, err := s.service.List(context.Background(), req)

	// check res
	if err != nil || res == nil || res.Total != 0 {
		s.passed = false
		mlog.Logger().Error(fmt.Sprintf("> test user list failed, res: %+v\n", res), zap.Error(err))
	}
}

func (s *userServiceTest) testCreate() {
	// make grpc req param
	req := &rpc_impl.User_CreateReq{
		OperatorId: s.testData.testID,
		UserName:   s.testData.testCreateUserName,
		Password:   utils.CalcSHA256(s.testData.testPassword),
		Permission: uint32(s.testData.testCreatePermission),
	}

	// invoke method
	_, err := s.service.Create(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName(s.testData.testCreateUserName)

	if err != nil || err2 != nil || newUser.Permission != s.testData.testCreatePermission {
		s.passed = false
		mlog.Logger().Error("> test user create failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}

	// for next steps: get 'user id' of create user
	s.testData.testCreateUserID = newUser.ID
}

func (s *userServiceTest) testLock() {
	// make grpc req param
	req := &rpc_impl.User_LockReq{
		OperatorId: s.testData.testID,
		UserId:     s.testData.testCreateUserID,
	}

	// invoke method
	_, err := s.service.Lock(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName(s.testData.testCreateUserName)

	if err != nil || err2 != nil || !newUser.IsLocked {
		s.passed = false
		mlog.Logger().Error("> test user lock failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func (s *userServiceTest) testUnlock() {
	// make grpc req param
	req := &rpc_impl.User_UnlockReq{
		OperatorId: s.testData.testID,
		UserId:     s.testData.testCreateUserID,
	}

	// invoke method
	_, err := s.service.Unlock(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName(s.testData.testCreateUserName)

	if err != nil || err2 != nil || newUser.IsLocked {
		s.passed = false
		mlog.Logger().Error("> test user unlock failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func (s *userServiceTest) testModifyInfo() {
	// make grpc req param
	req := &rpc_impl.User_ModifyInfoReq{
		OperatorId: s.testData.testCreateUserID,
		UserId:     s.testData.testCreateUserID,
		CurrPwd:    utils.CalcSHA256(s.testData.testPassword),
		Nickname:   s.testData.testModifyInfoUserName,
	}

	// invoke method
	_, err := s.service.ModifyInfo(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName(s.testData.testCreateUserName)

	if err != nil || err2 != nil || newUser.Nickname != s.testData.testModifyInfoUserName {
		s.passed = false
		mlog.Logger().Error("> test user modify info failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func (s *userServiceTest) testModifyPermission() {
	// make grpc req param
	req := &rpc_impl.User_ModifyPermissionReq{
		OperatorId: s.testData.testID,
		UserId:     s.testData.testCreateUserID,
		Permission: uint32(s.testData.testModifyPermissionP),
	}

	// invoke method
	_, err := s.service.ModifyPermission(context.Background(), req)

	// check res
	newUser, err2 := getUserByUserName(s.testData.testCreateUserName)

	if err != nil || err2 != nil || newUser.Permission != s.testData.testModifyPermissionP {
		s.passed = false
		mlog.Logger().Error("> test user modify permission failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func (s *userServiceTest) beforeTest(t *testing.T) {
	initialize.InitFromFile("server_test.json", mdb.Init, config.Init, db.Init)

	s.passed = true

	// create test table
	err := mtest.CreateTestTable_Postgresql([]interface{}{
		(*model.User)(nil),
	})
	if err != nil {
		t.Fail()
	}

	// prepare test data, include default table data
	s.testData.loadTestData()

	// create global service instance
	s.service = GetUserServer()
}

func (s *userServiceTest) afterTest(t *testing.T) {
	if s.passed {
		mtest.DropTestTable_Postgresql([]interface{}{
			(*model.User)(nil),
		})
	} else {
		t.Fail()
	}
}

func getUserByUserName(userName string) (*model.User, error) {
	user := &model.User{}
	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(user).Where(model.User_UserName+" = ?", userName).Select()
	})

	if err != nil {
		user = nil
	}

	return user, err
}
