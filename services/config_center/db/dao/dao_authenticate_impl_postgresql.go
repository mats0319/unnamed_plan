package dao

import (
    mdb "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
)

type AuthenticatePostgresql model.Administrator

var _ AuthenticateDao = (*AuthenticatePostgresql)(nil)

func (a *AuthenticatePostgresql) QueryOneByUserName(userName string) (*model.Administrator, error) {
    user := &model.Administrator{}

    err := mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        return conn.PostgresqlConn.Model(user).Where(model.Administrator_UserName + " = ?", userName).First()
    })
    if err != nil {
        user = nil
    }

    return user, err
}
