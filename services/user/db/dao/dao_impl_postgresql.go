package dao

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"time"
)

type UserPostgresql model.User

var _ UserDao = (*UserPostgresql)(nil)

func (u *UserPostgresql) Insert(user *model.User) (err error) {
	if len(user.ID) < 1 {
		user.Common = model.NewCommon()
	}

	return mdb.DB().WithTx(func(conn mdb.Conn) error {
		_, err = conn.PostgresqlConn.Model(user).Insert()
		return err
	})
}

func (u *UserPostgresql) Query(userIDs []string) (users []*model.User, err error) {
	err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
		return conn.PostgresqlConn.Model(&users).Where(model.Common_ID+" in (?)", pg.In(userIDs)).Select()
	})
	if err != nil {
		users = nil
	}

	return
}

func (u *UserPostgresql) QueryOne(userID string) (user *model.User, err error) {
	user = &model.User{}

	condition := fmt.Sprintf("%s = ? and %s = ?", model.User_IsLocked, model.Common_ID)

	err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
		return conn.PostgresqlConn.Model(user).Where(condition, false, userID).First()
	})
	if err != nil {
		user = nil
	}

	return
}

func (u *UserPostgresql) QueryOneByUserName(userName string) (user *model.User, err error) {
	user = &model.User{}

	condition := fmt.Sprintf("%s = ? and %s = ?", model.User_IsLocked, model.User_UserName)

	err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
		return conn.PostgresqlConn.Model(user).Where(condition, false, userName).First()
	})
	if err != nil {
		user = nil
	}

	return
}

// QueryPageLEPermission
/**
Core: sub-query
	select *
	from users u
	where "permission" <= (
		select "permission"
		from users u
		where user_id = 'user id'
	) and user_id != 'user id';
*/
func (u *UserPostgresql) QueryPageLEPermission(
	pageSize int,
	pageNum int,
	userID string,
) (users []*model.User, count int, err error) {
	err = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
		permission := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_Permission).
			Where(model.Common_ID+" = ?", userID)

		count, err = conn.PostgresqlConn.Model(&users).
			Where(model.User_Permission+" <= (?)", permission).
			Where(model.Common_ID+" != ?", userID).
			Order(model.User_Permission + " DESC").
			Offset((pageNum - 1) * pageSize).Limit(pageSize).SelectAndCount()

		return err
	})
	if err != nil {
		users = nil
		count = 0
	}

	return
}

func (u *UserPostgresql) UpdateColumnsByUserID(user *model.User, columns ...string) error {
	user.UpdateTime = time.Duration(time.Now().Unix())

	return mdb.DB().WithTx(func(conn mdb.Conn) error {
		query := conn.PostgresqlConn.Model(user).Column(model.Common_UpdateTime)
		for i := range columns {
			query.Column(columns[i])
		}

		_, err := query.Where(model.Common_ID + " = ?" + model.Common_ID).Update()
		return err
	})
}
