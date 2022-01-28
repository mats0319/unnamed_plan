package dao

import "github.com/mats9693/unnamed_plan/services/shared/db/model"

type TaskDao interface {
    Insert(task *model.Task) error

    QueryByPoster(userID string) (tasks []*model.Task, count int, err error)

    QueryOne(taskID string) (*model.Task, error)

    UpdateColumnsByTaskID(task *model.Task, columns ...string) error
}
