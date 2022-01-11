package dao

import (
    "fmt"
    "github.com/mats9693/unnamed_plan/services/shared/const"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
    "github.com/mats9693/unnamed_plan/services/shared/utils"
    "github.com/mats9693/utils/uuid"
    "time"
)

type ThinkNotePostgresql model.ThinkingNote

var _ ThinkingNoteDao = (*ThinkNotePostgresql)(nil)

func (tn *ThinkNotePostgresql) Insert(data *model.ThinkingNote) error {
    if len(data.ID) < 1 {
        data.Common = model.NewCommon()
    }

    if len(data.NoteID) < 1 {
        data.NoteID = uuid.New()
    }

    return mdb.DB().WithTx(func(conn mdb.Conn) error {
        _, err := conn.PostgresqlConn.Model(data).Insert()
        return err
    })
}

func (tn *ThinkNotePostgresql) QueryOne(thinkingNoteID string) (note *model.ThinkingNote, err error) {
    note = &model.ThinkingNote{}

    condition := fmt.Sprintf("%s = ? and %s = ?", model.ThinkingNote_IsDeleted, model.ThinkingNote_NoteID)

    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        return conn.PostgresqlConn.Model(note).Where(condition, false, thinkingNoteID).Select()
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
        count, err = conn.PostgresqlConn.Model(&notes).Where(model.ThinkingNote_IsDeleted+" = ?", false).
            Where(model.ThinkingNote_WriteBy+" = ?", userID).
            Order(model.Common_UpdateTime + " DESC").
            Offset((pageNum - 1) * pageSize).Limit(pageSize).SelectAndCount()

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
		select "user_id"
		from users u
		where "permission" <= (
			select "permission"
			from users u2
			where user_id = 'user id'
		)
	);
*/
func (tn *ThinkNotePostgresql) QueryPageInPublic(
    pageSize int,
    pageNum int,
    userID string,
) (notes []*model.ThinkingNote, count int, err error) {
    err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        permission := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.User_UserID+" = ?", userID)
        userIDs := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_UserID).Where(model.User_Permission+" <= (?)", permission)

        count, err = conn.PostgresqlConn.Model(&notes).
            Where(model.ThinkingNote_IsDeleted+" = ?", false).
            Where(model.ThinkingNote_IsPublic+" = ?", true).
            Where(model.ThinkingNote_WriteBy+" in (?)", userIDs).
            Order(model.Common_UpdateTime + " DESC").
            Offset((pageNum - 1) * pageSize).Limit(pageSize).SelectAndCount()

        return err
    })
    if err != nil {
        notes = nil
        count = 0
    }

    return
}

func (tn *ThinkNotePostgresql) UpdateColumnsByNoteID(thinkingNote *model.ThinkingNote, columns ...string) error {
    thinkingNote.UpdateTime = time.Duration(time.Now().Unix())

    return mdb.DB().WithTx(func(conn mdb.Conn) error {
        query := conn.PostgresqlConn.Model(thinkingNote).Column(model.Common_UpdateTime)
        for i := range columns {
            query.Column(columns[i])
        }

        res, err := query.Where(model.ThinkingNote_IsDeleted+" = ?", false).
            Where(model.ThinkingNote_NoteID + " = ?" + model.ThinkingNote_NoteID).Update()
        if err != nil {
            return err
        }

        if res.RowsAffected() < 0 {
            return utils.NewError(mconst.Error_NoteAlreadyDeleted)
        }

        return nil
    })
}
