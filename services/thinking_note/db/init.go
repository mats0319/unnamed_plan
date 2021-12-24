package db

import (
    "github.com/mats9693/unnamed_plan/services/thinking_note/db/dao"
    "github.com/mats9693/utils/toy_server/const"
    "github.com/mats9693/utils/toy_server/db"
    "github.com/mats9693/utils/toy_server/log"
    "os"
)

var thinkingNoteDaoIns dao.ThinkingNoteDao

func GetThinkingNoteDao() dao.ThinkingNoteDao {
    return thinkingNoteDaoIns
}

func init() {
    switch mdb.DB().GetDB() {
    case mconst.DB_PostgreSQL:
        thinkingNoteDaoIns = &dao.ThinkNotePostgresql{}
    default:
        mlog.Logger().Error(mconst.Error_UnsupportedDB)
        os.Exit(-1)
    }

    mlog.Logger().Info("> Database instance init finish.")
}
