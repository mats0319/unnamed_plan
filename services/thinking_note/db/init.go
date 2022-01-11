package db

import (
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "github.com/mats9693/unnamed_plan/services/thinking_note/db/dao"
    "os"
)

var thinkingNoteDaoIns dao.ThinkingNoteDao

func GetThinkingNoteDao() dao.ThinkingNoteDao {
    return thinkingNoteDaoIns
}

func init() {
    switch mdb.DB().GetDBMSName() {
    case mconst.DB_PostgreSQL:
        thinkingNoteDaoIns = &dao.ThinkNotePostgresql{}
    default:
        mlog.Logger().Error(mconst.Error_UnsupportedDB)
        os.Exit(-1)
    }

    mlog.Logger().Info("> Database instance init finish.")
}
