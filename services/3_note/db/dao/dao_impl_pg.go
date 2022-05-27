package dao

import (
	"github.com/mats9693/unnamed_plan/services/shared/const"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"github.com/mats9693/unnamed_plan/services/shared/utils"
	"time"
)

type NotePostgresql model.Note

var _ NoteDao = (*NotePostgresql)(nil)

func (tn *NotePostgresql) Insert(data *model.Note) error {
	if len(data.ID) < 1 {
		data.Common = model.NewCommon()
	}

	return mdb.DB().WithTx(func(conn mdal.Conn) error {
		_, err := conn.PostgresqlConn.Model(data).Insert()
		return err
	})
}

func (tn *NotePostgresql) QueryOne(noteID string) (note *model.Note, err error) {
	note = &model.Note{}

	err = mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(note).
			Where(model.Note_IsDeleted+" = ?", false).
			Where(model.Common_ID+" = ?", noteID).
			Select()
	})
	if err != nil {
		note = nil
	}

	return
}

func (tn *NotePostgresql) QueryPageByWriter(
	pageSize int,
	pageNum int,
	userID string,
) (notes []*model.Note, count int, err error) {
	err = mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		count, err = conn.PostgresqlConn.Model(&notes).
			Where(model.Note_IsDeleted+" = ?", false).
			Where(model.Note_WriteBy+" = ?", userID).
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
	from note tn
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
func (tn *NotePostgresql) QueryPageInPublic(
	pageSize int,
	pageNum int,
	userID string,
) (notes []*model.Note, count int, err error) {
	err = mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		permission := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.Common_ID+" = ?", userID)
		userIDs := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.Common_ID).Where(model.User_Permission+" <= (?)", permission)

		count, err = conn.PostgresqlConn.Model(&notes).
			Where(model.Note_IsDeleted+" = ?", false).
			Where(model.Note_IsPublic+" = ?", true).
			Where(model.Note_WriteBy+" in (?)", userIDs).
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

func (tn *NotePostgresql) UpdateColumnsByNoteID(note *model.Note, columns ...string) error {
	note.UpdateTime = time.Duration(time.Now().Unix())
	note.OptimisticLock++

	return mdb.DB().WithTx(func(conn mdal.Conn) error {
		query := conn.PostgresqlConn.Model(note).Column(model.Common_UpdateTime, model.Common_OptimisticLock)
		for i := range columns {
			query.Column(columns[i])
		}

		res, err := query.Where(model.Note_IsDeleted+" = ?", false).
			Where(model.Common_ID+" = ?"+model.Common_ID).
			Where(model.Common_OptimisticLock+" = ?", note.OptimisticLock-1).Update()
		if err != nil {
			return err
		}

		if res.RowsAffected() <= 0 {
			return utils.NewError(mconst.Error_NoAffectedRows)
		}

		return nil
	})
}
