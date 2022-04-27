package rpc

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/init"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/impl"
	"github.com/mats9693/unnamed_plan/services/shared/proto/mock"
	"github.com/mats9693/unnamed_plan/services/shared/test"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/mats9693/unnamed_plan/services/task/config"
	"github.com/mats9693/unnamed_plan/services/task/db"
	"go.uber.org/zap"
	"testing"
)

type taskServiceTest struct {
	service *taskServerImpl
	passed  bool

	testData *testData
}

func TestTaskService(t *testing.T) {
	serviceIns := &taskServiceTest{testData: &testData{}}

	serviceIns.beforeTest(t)

	serviceIns.testList()
	serviceIns.testCreate()
	serviceIns.testModify()

	serviceIns.afterTest(t)
}

func (s *taskServiceTest) testList() {
	// make grpc req param
	req := &rpc_impl.Task_ListReq{
		OperatorId: s.testData.testUserID,
	}

	// invoke method
	res, err := s.service.List(context.Background(), req)

	// check res
	if err != nil || res == nil || res.Total != 1 || res.Tasks[0].TaskName != s.testData.testTaskName {
		s.passed = false
		mlog.Logger().Error(fmt.Sprintf("> test task list failed, res: %+v\n", res), zap.Error(err))
	}
}

func (s *taskServiceTest) testCreate() {
	// make grpc req param
	req := &rpc_impl.Task_CreateReq{
		OperatorId:  s.testData.testUserID,
		TaskName:    s.testData.testCreatedTaskName,
		Description: s.testData.testCreatedTaskDesp,
		PreTaskIds:  nil,
	}

	// invoke method
	_, err := s.service.Create(context.Background(), req)

	// check res
	task, err2 := getTaskByTaskName(s.testData.testCreatedTaskName)

	if err != nil || err2 != nil || task.Description != s.testData.testCreatedTaskDesp {
		s.passed = false
		mlog.Logger().Error("> test task create failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}

	// for next steps
	s.testData.testCreatedTaskID = task.ID
}

func (s *taskServiceTest) testModify() {
	// make grpc req param
	req := &rpc_impl.Task_ModifyReq{
		OperatorId:  s.testData.testUserID,
		TaskId:      s.testData.testCreatedTaskID,
		Password:    utils.CalcSHA256(s.testData.testPassword),
		TaskName:    s.testData.testCreatedTaskName,
		Description: s.testData.testCreatedTaskDesp,
		PreTaskIds:  nil,
		Status:      uint32(s.testData.testModifyTaskStatus), // modify
	}

	// invoke method
	_, err := s.service.Modify(context.Background(), req)

	// check res
	note, err2 := getTaskByTaskName(s.testData.testCreatedTaskName)

	if err != nil || err2 != nil || note.Status != s.testData.testModifyTaskStatus {
		s.passed = false
		mlog.Logger().Error("> test task modify failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func (s *taskServiceTest) beforeTest(t *testing.T) {
	initialize.InitFromFile("../config_test.json", mdb.Init, config.Init, db.Init)

	s.passed = true

	// create test table
	err := mtest.CreateTestTable_Postgresql([]interface{}{
		(*model.User)(nil),
		(*model.Task)(nil),
	})
	if err != nil {
		t.Fail()
	}

	// prepare test data, include default table data
	s.testData.loadTestData()

	// create global service instance
	s.service = taskServerImplIns

	userClientMock := mock_rpc_impl.NewMockIUserClient(gomock.NewController(t))
	userClientMock.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	s.service.userClient = userClientMock
}

func (s *taskServiceTest) afterTest(t *testing.T) {
	if s.passed {
		mtest.DropTestTable_Postgresql([]interface{}{
			(*model.User)(nil),
			(*model.Task)(nil),
		})
	} else {
		t.Fail()
	}
}

func getTaskByTaskName(taskName string) (*model.Task, error) {
	task := &model.Task{}
	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(task).Where(model.Task_TaskName+" = ?", taskName).Select()
	})

	if err != nil {
		task = nil
	}

	return task, err
}
