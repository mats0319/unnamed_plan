package dao

import "github.com/mats9693/unnamed_plan/services/shared/db/model"

type ThinkingNoteDao interface {
    // Insert insert one note
    Insert(*model.ThinkingNote) error

    // QueryOne query one undeleted note by 'note id'
    QueryOne(thinkingNoteID string) (*model.ThinkingNote, error)

    // QueryPageByWriter query notes that write by designated user(designate by user id)
    // result not contains notes marked as 'deleted', result order by 'update time' desc
    QueryPageByWriter(pageSize int, pageNum int, userID string) (notes []*model.ThinkingNote, count int, err error)

    // QueryPageInPublic query notes satisfy following requirements:
    //   1. note is 'public'
    //   2. permission of writer is less than or equal to designated user(designate by user id)
    // result not contains notes marked as 'deleted', result order by 'update time' desc
    QueryPageInPublic(pageSize int, pageNum int, userID string) (notes []*model.ThinkingNote, count int, err error)

    // UpdateColumnsByNoteID update designated 'columns' on designated undeleted note(designate by note id)
    UpdateColumnsByNoteID(data *model.ThinkingNote, columns ...string) error
}
