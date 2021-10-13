package dao

import (
	"github.com/go-pg/pg/v10/orm"
	"time"

	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	mdb "github.com/mats9693/utils/toy_server/db"
	"github.com/mats9693/utils/uuid"
)

type ThinkNote model.ThinkingNote

var thinkingNoteIns = &ThinkNote{}

func GetThinkingNote() *ThinkNote {
	return thinkingNoteIns
}

// QueryFirst 仅选择未删除的笔记
func (tn *ThinkNote) QueryFirst(condition string, params ...interface{}) (note *model.ThinkingNote, err error) {
	note = &model.ThinkingNote{}

	err = mdb.WithNoTx(func(conn orm.DB) error {
		return conn.Model(note).Where(model.ThinkingNote_IsDeleted+" = ?", false).Where(condition, params...).Select()
	})
	if err != nil {
		note = nil
	}

	return
}

// QueryPageByWriter 不查询已删除的记录，按照更新时间降序
func (tn *ThinkNote) QueryPageByWriter(
	pageSize int,
	pageNum int,
	userID string,
) (notes []*model.ThinkingNote, count int, err error) {
	err = mdb.WithNoTx(func(conn orm.DB) error {
		count, err = conn.Model(&notes).Where(model.ThinkingNote_IsDeleted+" = ?", false).
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

// QueryPageInPublic 获取公开笔记列表，要求编辑者权限等级不高于指定用户（通过userID指定），分页，不查询已删除的记录，按照更新时间降序
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
func (tn *ThinkNote) QueryPageInPublic(
	pageSize int,
	pageNum int,
	userID string,
) (notes []*model.ThinkingNote, count int, err error) {
	err = mdb.WithNoTx(func(conn orm.DB) error {
		permission := conn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.User_UserID+" = ?", userID)
		userIDs := conn.Model((*model.User)(nil)).Column(model.User_UserID).Where(model.User_Permission+" <= (?)", permission)

		count, err = conn.Model(&notes).Where(model.ThinkingNote_IsDeleted+" = ?", false).
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

// UpdateColumnsByNoteID 仅可操作未删除的笔记
func (tn *ThinkNote) UpdateColumnsByNoteID(data *model.ThinkingNote, columns ...string) (err error) {
	data.UpdateTime = time.Duration(time.Now().Unix())

	return mdb.WithTx(func(conn orm.DB) error {
		query := conn.Model(data).Column(model.Common_UpdateTime)
		for i := range columns {
			query.Column(columns[i])
		}

		_, err = query.Where(model.ThinkingNote_IsDeleted+" = ?", false).
			Where(model.ThinkingNote_NoteID + " = ?" + model.ThinkingNote_NoteID).Update()

		return err
	})
}
