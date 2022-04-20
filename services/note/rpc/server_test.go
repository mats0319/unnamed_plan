package rpc

import (
    "context"
    "fmt"
    "github.com/golang/mock/gomock"
    "github.com/mats9693/unnamed_plan/services/note/db"
    "github.com/mats9693/unnamed_plan/services/shared/db"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
    "github.com/mats9693/unnamed_plan/services/shared/init"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/shared/proto/impl"
    "github.com/mats9693/unnamed_plan/services/shared/proto/mock"
    "github.com/mats9693/unnamed_plan/services/shared/test"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "go.uber.org/zap"
    "testing"
)

type noteServiceTest struct {
    service *noteServerImpl
    passed bool

    testData *testData
}

func TestNoteService(t *testing.T) {
    serviceIns := &noteServiceTest{testData: &testData{}}

    serviceIns.beforeTest(t)

    serviceIns.testListByWriter()
    serviceIns.testListPublic()
    serviceIns.testCreate()
    serviceIns.testModify()
    serviceIns.testDelete()

    serviceIns.afterTest(t)
}

func (s *noteServiceTest) testListByWriter() {
    // make grpc req param
    req := &rpc_impl.Note_ListByWriterReq{
        OperatorId: s.testData.testUserID,
        Page: &rpc_impl.Pagination{
            PageSize: 10,
            PageNum:  1,
        },
    }

    // invoke method
    res, err := s.service.ListByWriter(context.Background(), req)

    // check res
    if err != nil || res == nil || res.Total != 1 || res.Notes[0].Content != s.testData.testNoteContent {
        s.passed = false
        mlog.Logger().Error(fmt.Sprintf("> test note list by writer failed, res: %+v\n", res), zap.Error(err))
    }
}

func (s *noteServiceTest) testListPublic() {
    // make grpc req param
    req := &rpc_impl.Note_ListPublicReq{
        OperatorId: s.testData.testUserID,
        Page: &rpc_impl.Pagination{
            PageSize: 10,
            PageNum:  1,
        },
    }

    // invoke method
    res, err := s.service.ListPublic(context.Background(), req)

    // check res
    if err != nil || res == nil || res.Total != 1 || res.Notes[0].Content != s.testData.testNoteContent {
        s.passed = false
        mlog.Logger().Error(fmt.Sprintf("> test note list public failed, res: %+v\n", res), zap.Error(err))
    }
}

func (s *noteServiceTest) testCreate() {
    // make grpc req param
    req := &rpc_impl.Note_CreateReq{
        OperatorId: s.testData.testUserID,
        Topic:      s.testData.testCreatedNoteTopic,
        Content:    s.testData.testCreatedNoteContent,
        IsPublic:   false,
    }

    // invoke method
    _, err := s.service.Create(context.Background(), req)

    // check res
    note, err2 := getNoteByTopic(s.testData.testCreatedNoteTopic)

    if err != nil || err2 != nil || note.Content != s.testData.testCreatedNoteContent {
        s.passed = false
        mlog.Logger().Error("> test note create failed",
            zap.NamedError("err", err),
            zap.NamedError("err2", err2))
    }

    // for next steps
    s.testData.testCreatedNoteID = note.ID
}

func (s *noteServiceTest) testModify() {
    // make grpc req param
    req := &rpc_impl.Note_ModifyReq{
        OperatorId: s.testData.testUserID,
        NoteId:     s.testData.testCreatedNoteID,
        Password:   utils.CalcSHA256(s.testData.testPassword),
        Topic:      s.testData.testCreatedNoteTopic,
        Content:    s.testData.testCreatedNoteContent,
        IsPublic:   true, // modify
    }

    // invoke method
    _, err := s.service.Modify(context.Background(), req)

    // check res
    note, err2 := getNoteByTopic(s.testData.testCreatedNoteTopic)

    if err != nil || err2 != nil || !note.IsPublic {
        s.passed = false
        mlog.Logger().Error("> test note modify failed",
            zap.NamedError("err", err),
            zap.NamedError("err2", err2))
    }
}

func (s *noteServiceTest) testDelete() {
    // make grpc req param
    req := &rpc_impl.Note_DeleteReq{
        OperatorId: s.testData.testUserID,
        NoteId:     s.testData.testCreatedNoteID,
        Password:   utils.CalcSHA256(s.testData.testPassword),
    }

    // invoke method
    _, err := s.service.Delete(context.Background(), req)

    // check res
    note, err2 := getNoteByTopic(s.testData.testCreatedNoteTopic)

    if err != nil || err2 != nil || !note.IsDeleted {
        s.passed = false
        mlog.Logger().Error("> test note delete failed",
            zap.NamedError("err", err),
            zap.NamedError("err2", err2))
    }
}

func (s *noteServiceTest) beforeTest(t *testing.T) {
    initialize.Init("../config_test.json", mdb.Init, db.Init)

    s.passed = true

    // create test table
    err := mtest.CreateTestTable_Postgresql([]interface{}{
        (*model.User)(nil),
        (*model.Note)(nil),
    })
    if err != nil {
        t.Fail()
    }

    // prepare test data, include default table data
    s.testData.loadTestData()

    // create global service instance
    s.service = noteServerImplIns

    userClientMock := mock_rpc_impl.NewMockIUserClient(gomock.NewController(t))
    userClientMock.EXPECT().Authenticate(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

    s.service.userClient = userClientMock
}

func (s *noteServiceTest) afterTest(t *testing.T) {
    if s.passed {
        mtest.DropTestTable_Postgresql([]interface{}{
            (*model.User)(nil),
            (*model.Note)(nil),
        })
    } else {
        t.Fail()
    }
}

func getNoteByTopic(topic string) (*model.Note, error) {
    note := &model.Note{}
    err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
        return conn.PostgresqlConn.Model(note).Where(model.Note_Topic+" = ?", topic).Select()
    })

    if err != nil {
        note = nil
    }

    return note, err
}
