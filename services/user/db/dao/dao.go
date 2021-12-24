package dao

import "github.com/mats9693/unnamed_plan/services/shared/db/model"

type UserDao interface {
    // Insert insert one user
    Insert(*model.User) error

    // Query query users by ids
    Query(userIDs []string) ([]*model.User, error)

    // QueryOne query one unlocked user by id
    QueryOne(userID string) (*model.User, error)

    // QueryOneByUserName query one unlocked user by 'user name'
    QueryOneByUserName(string) (*model.User, error)

    // QueryPageLEPermission query users that permission less than or equal to designate user(designate by id)
    // result not contains designate user, result order by 'permission' desc
    QueryPageLEPermission(pageSize int, pageNum int, userID string) (users []*model.User, count int, err error)

    // UpdateColumnsByUserID update designated 'columns' on designated user(designate by id)
    UpdateColumnsByUserID(user *model.User, columns ...string) error
}
