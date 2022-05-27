package rpc

import (
	"github.com/mats9693/unnamed_plan/services/3_note/db"
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

	testNotes []*model.Note

	testNoteTopic   string
	testNoteContent string

	testCreatedNoteTopic   string
	testCreatedNoteContent string
	testCreatedNoteID      string // var
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

	td.testNoteTopic = "test topic"
	td.testNoteContent = "test content"

	td.testNotes = []*model.Note{
		{
			WriteBy:   td.testUserID,
			Topic:     td.testNoteTopic,
			Content:   td.testNoteContent,
			IsPublic:  true,
			IsDeleted: false,
			Common:    model.Common{},
		},
	}

	setTestNote(td.testNotes)

	td.testCreatedNoteTopic = "test create note topic"
	td.testCreatedNoteContent = "test create note content"
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

func setTestNote(notes []*model.Note) {
	for i := range notes {
		err := db.GetNoteDao().Insert(notes[i])
		if err != nil {
			mlog.Logger().Error("set default note failed", zap.Int("index", i), zap.Error(err))
			os.Exit(-1)
		}
	}
}
