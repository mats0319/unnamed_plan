package rpc

import (
	"github.com/mats9693/unnamed_plan/services/cloud_file/db"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
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

	testFiles []*model.CloudFile

	testFileID        string
	testFileName      string
	testExtensionName string

	testUploadFileContent string
	testUploadFileName    string

	testCreatedFileID string // var
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

	td.testFileID = "test file id"
	td.testFileName = "test file name"
	td.testExtensionName = "test extension name"

	td.testFiles = []*model.CloudFile{
		{
			FileID:           td.testFileID,
			UploadedBy:       td.testUserID,
			FileName:         td.testFileName,
			ExtensionName:    td.testExtensionName,
			LastModifiedTime: 1000,
			FileSize:         1000,
			IsPublic:         true,
			IsDeleted:        false,
		},
	}

	setTestCloudFile(td.testFiles)

    td.testUploadFileContent = "test upload file content"
    td.testUploadFileName = "test upload file name"
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

func setTestCloudFile(files []*model.CloudFile) {
	for i := range files {
		err := db.GetCloudFileDao().Insert(files[i])
		if err != nil {
			mlog.Logger().Error("set default cloud file failed", zap.Int("index", i), zap.Error(err))
			os.Exit(-1)
		}
	}
}
