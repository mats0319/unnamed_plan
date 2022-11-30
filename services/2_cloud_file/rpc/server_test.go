package rpc

import (
	"context"
	"fmt"
	"github.com/mats9693/unnamed_plan/services/2_cloud_file/config"
	"github.com/mats9693/unnamed_plan/services/2_cloud_file/db"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/init"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"github.com/mats9693/unnamed_plan/services/shared/proto/go"
	"github.com/mats9693/unnamed_plan/services/shared/proto/mock"
	"github.com/mats9693/unnamed_plan/services/shared/test"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"go.uber.org/zap"
	"testing"
)

type cloudFileServiceTest struct {
	service *cloudFileServerImpl
	passed  bool

	testData *testData
}

func TestCloudFileService(t *testing.T) {
	serviceIns := &cloudFileServiceTest{testData: &testData{}}

	serviceIns.beforeTest(t)

	serviceIns.testList()
	serviceIns.testUpload()
	serviceIns.testModify()
	serviceIns.testDelete()

	serviceIns.afterTest(t)
}

func (s *cloudFileServiceTest) testList() {
	// make grpc req param
	req := &rpc_impl.CloudFile_ListReq{
		Rule:       1,
		OperatorId: s.testData.testUserID,
		Page: &rpc_impl.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
	}

	// invoke method
	res, err := s.service.List(context.Background(), req)

	// check res
	if err != nil || res == nil || res.Err != nil || res.Total != 1 || res.Files[0].FileId != s.testData.testFileID {
		s.passed = false
		mlog.Logger().Error(fmt.Sprintf("> test list cloud file failed, res: %+v\n", res), zap.Error(err))
	}
}

func (s *cloudFileServiceTest) testUpload() {
	// make grpc req param
	req := &rpc_impl.CloudFile_UploadReq{
		OperatorId:       s.testData.testUserID,
		File:             []byte(s.testData.testUploadFileContent),
		FileName:         s.testData.testUploadFileName,
		ExtensionName:    s.testData.testExtensionName,
		FileSize:         int64(len(s.testData.testUploadFileContent)),
		LastModifiedTime: 1000,
		IsPublic:         false,
	}

	// invoke method
	res, err := s.service.Upload(context.Background(), req)

	// check res
	file, err2 := getFileByFileName(s.testData.testUploadFileName)

	if err != nil || err2 != nil || res.Err != nil || file.UploadedBy != s.testData.testUserID {
		s.passed = false
		mlog.Logger().Error("> test cloud file upload failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}

	// for next steps
	s.testData.testCreatedFileID = file.FileID
}

func (s *cloudFileServiceTest) testModify() {
	// make grpc req param
	req := &rpc_impl.CloudFile_ModifyReq{
		OperatorId:       s.testData.testUserID,
		FileId:           s.testData.testCreatedFileID,
		Password:         utils.CalcSHA256(s.testData.testPassword),
		FileName:         s.testData.testUploadFileName,
		ExtensionName:    s.testData.testExtensionName,
		IsPublic:         true, // modify
		File:             []byte(s.testData.testUploadFileContent),
		FileSize:         int64(len(s.testData.testUploadFileContent)),
		LastModifiedTime: 1000,
	}

	// invoke method
	res, err := s.service.Modify(context.Background(), req)

	// check res
	file, err2 := getFileByFileName(s.testData.testUploadFileName)

	if err != nil || err2 != nil || res.Err != nil || !file.IsPublic {
		s.passed = false
		mlog.Logger().Error("> test cloud file modify failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func (s *cloudFileServiceTest) testDelete() {
	// make grpc req param
	req := &rpc_impl.CloudFile_DeleteReq{
		OperatorId: s.testData.testUserID,
		FileId:     s.testData.testCreatedFileID,
		Password:   utils.CalcSHA256(s.testData.testPassword),
	}

	// invoke method
	res, err := s.service.Delete(context.Background(), req)

	// check res
	file, err2 := getFileByFileName(s.testData.testUploadFileName)

	if err != nil || err2 != nil || res.Err != nil || !file.IsDeleted {
		s.passed = false
		mlog.Logger().Error("> test cloud file delete failed",
			zap.NamedError("err", err),
			zap.NamedError("err2", err2))
	}
}

func (s *cloudFileServiceTest) beforeTest(t *testing.T) {
	err := initialize.InitFromFile("server_test.json", mdb.Init, config.Init, db.Init)
	if err != nil {
		t.Fail()
	}

	s.passed = true

	// create test table
	err = mtest.CreateTestTable_Postgresql([]interface{}{
		(*model.User)(nil),
		(*model.CloudFile)(nil),
	})
	if err != nil {
		t.Fail()
	}

	// prepare test data, include default table data
	s.testData.loadTestData()

	// create global service instance
	s.service = cloudFileServerImplIns

	mock_rpc_impl.MockUserAuthenticate()
}

func (s *cloudFileServiceTest) afterTest(t *testing.T) {
	if s.passed {
		mtest.DropTestTable_Postgresql([]interface{}{
			(*model.User)(nil),
			(*model.CloudFile)(nil),
		})
	} else {
		t.Fail()
	}
}

func getFileByFileName(fileName string) (*model.CloudFile, error) {
	file := &model.CloudFile{}
	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(file).Where(model.CloudFile_FileName+" = ?", fileName).Select()
	})

	if err != nil {
		file = nil
	}

	return file, err
}
