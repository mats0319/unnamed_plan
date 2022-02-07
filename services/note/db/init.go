package db

import (
	"github.com/mats9693/unnamed_plan/services/note/db/dao"
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"os"
)

var noteDaoIns dao.NoteDao

func GetNoteDao() dao.NoteDao {
	return noteDaoIns
}

func init() {
	switch mdb.DB().GetDBMSName() {
	case mconst.DB_PostgreSQL:
		noteDaoIns = &dao.NotePostgresql{}
	default:
		mlog.Logger().Error(mconst.Error_UnsupportedDB)
		os.Exit(-1)
	}

	mlog.Logger().Info("> Database instance init finish.")
}
