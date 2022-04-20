package rpc

import (
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"github.com/mats9693/unnamed_plan/services/task/db"
	"go.uber.org/zap"
	"os"
)

type testData struct {
	testUsers []*model.User

	testUserName   string
	testNickname   string
	testPassword   string
	testPermission uint8
	testUserID     string

	testTasks []*model.Task

	testTaskName string
	testTaskDesp string

	testCreatedTaskName string
	testCreatedTaskDesp string
	testCreatedTaskID   string // var

	testModifyTaskStatus uint8
}

func (td *testData) loadTestData() {
	td.testUserName = "test user"
	td.testNickname = "test nickname"
	td.testPassword = "test password"
	td.testPermission = 8
	td.testUserID = "test id"

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
				ID: td.testUserID,
			},
		},
	}

	setTestUser(td.testUsers)

	td.testTaskName = "test task name"
	td.testTaskDesp = "test task desp"

	td.testTasks = []*model.Task{
		{
			TaskName:    td.testTaskName,
			PostedBy:    td.testUserID,
			Description: td.testTaskDesp,
			PreTaskIDs:  nil,
			Status:      1,
			Common:      model.Common{},
		},
	}

	setTestTask(td.testTasks)

	td.testCreatedTaskName = "test created task name"
	td.testCreatedTaskDesp = "test created task desp"

	td.testModifyTaskStatus = 8
}

func setTestUser(users []*model.User) {
	for i := range users {
		users[i].Password = utils.CalcSHA256(users[i].Password)
		users[i].Password = utils.CalcSHA256(users[i].Password, users[i].Salt)

		if len(users[i].ID) < 1 {
			users[i].Common = model.NewCommon()
		}

		err := mdb.DB().WithTx(func(conn mdal.Conn) error {
			_, err := conn.PostgresqlConn.Model(users[i]).Insert()
			return err
		})
		if err != nil {
			mlog.Logger().Error("set default users failed", zap.Int("index", i), zap.Error(err))
			os.Exit(-1)
		}
	}
}

func setTestTask(tasks []*model.Task) {
	for i := range tasks {
		err := db.GetTaskDao().Insert(tasks[i])
		if err != nil {
			mlog.Logger().Error("set default task failed", zap.Int("index", i), zap.Error(err))
			os.Exit(-1)
		}
	}
}
