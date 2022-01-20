package dao

import (
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "time"
)

type ThinkNotePostgresql model.ThinkingNote

var _ ThinkingNoteDao = (*ThinkNotePostgresql)(nil)

func (tn *ThinkNotePostgresql) Insert(data *model.ThinkingNote) error {
    if len(data.ID) < 1 {
        data.Common = model.NewCommon()
    }

    return mdb.DB().WithTx(func(conn mdb.Conn) error {
        _, err := conn.PostgresqlConn.Model(data).Insert()
        return err
    })
}

func (tn *ThinkNotePostgresql) QueryOne(noteID string) (note *model.ThinkingNote, err error) {
    note = &model.ThinkingNote{}

    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        return conn.PostgresqlConn.Model(note).
            Where(model.ThinkingNote_IsDeleted + " = ?", false).
            Where(model.Common_ID + " = ?", noteID).
            Select()
    })
    if err != nil {
        note = nil
    }

    return
}

func (tn *ThinkNotePostgresql) QueryPageByWriter(
    pageSize int,
    pageNum int,
    userID string,
) (notes []*model.ThinkingNote, count int, err error) {
    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        count, err = conn.PostgresqlConn.Model(&notes).
            Where(model.ThinkingNote_IsDeleted+" = ?", false).
            Where(model.ThinkingNote_WriteBy+" = ?", userID).
            Order(model.Common_UpdateTime + " DESC").
            Offset((pageNum - 1) * pageSize).
            Limit(pageSize).
            SelectAndCount()

        return err
    })
    if err != nil {
        notes = nil
        count = 0
    }

    return
}

// QueryPageInPublic
/**
Core: sub-query
	select *
	from thinking_note tn
	where tn.write_by in (
		select "id"
		from users u
		where "permission" <= (
			select "permission"
			from users u2
			where id = 'user id'
		)
	);
*/
func (tn *ThinkNotePostgresql) QueryPageInPublic(
    pageSize int,
    pageNum int,
    userID string,
) (notes []*model.ThinkingNote, count int, err error) {
    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        permission := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.Common_ID+" = ?", userID)
        userIDs := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.Common_ID).Where(model.User_Permission+" <= (?)", permission)

        count, err = conn.PostgresqlConn.Model(&notes).
            Where(model.ThinkingNote_IsDeleted+" = ?", false).
            Where(model.ThinkingNote_IsPublic+" = ?", true).
            Where(model.ThinkingNote_WriteBy+" in (?)", userIDs).
            Order(model.Common_UpdateTime + " DESC").
            Offset((pageNum - 1) * pageSize).
            Limit(pageSize).
            SelectAndCount()

        return err
    })
    if err != nil {
        notes = nil
        count = 0
    }

    return
}

func (tn *ThinkNotePostgresql) UpdateColumnsByNoteID(note *model.ThinkingNote, columns ...string) error {
    note.UpdateTime = time.Duration(time.Now().Unix())

    return mdb.DB().WithTx(func(conn mdb.Conn) error {
        query := conn.PostgresqlConn.Model(note).Column(model.Common_UpdateTime)
        for i := range columns {
            query.Column(columns[i])
        }

        res, err := query.Where(model.ThinkingNote_IsDeleted+" = ?", false).
            Where(model.Common_ID + " = ?" + model.Common_ID).Update()
        if err != nil {
            return err
        }

        if res.RowsAffected() < 0 {
            return utils.NewError(mconst.Error_NoteAlreadyDeleted)
        }

        return nil
    })
}
