package dao

import (
    "github.com/go-pg/pg/v10/orm"

    "github.com/mats9693/unnamed_plan/admin_data/db/model"
    mdb "github.com/mats9693/utils/toy_server/db"
    "github.com/mats9693/utils/uuid"
)

type ThinkNote model.ThinkingNote

var thinkingNoteIns = &ThinkNote{}

func GetThinkingNote() *ThinkNote {
    return thinkingNoteIns
}

func (tn *ThinkNote) Insert(data *model.ThinkingNote) error {
    if len(data.ID) < 1 {
        data.Common = model.NewCommon()
    }

    if len(data.NoteID) < 1 {
        data.NoteID = uuid.New()
    }

    return mdb.WithTx(func(conn orm.DB) error {
        _, err := conn.Model(data).Insert()
        return err
    })
}
