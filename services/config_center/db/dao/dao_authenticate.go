package dao

import "github.com/mats9693/unnamed_plan/services/shared/db/model"

type AuthenticateDao interface {
    // QueryOneByUserName query one user by 'user name'
    QueryOneByUserName(string) (*model.Administrator, error)
}
