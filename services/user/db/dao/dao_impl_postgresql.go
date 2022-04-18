package dao

import (
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
	"time"
)

type UserPostgresql model.User

var _ UserDao = (*UserPostgresql)(nil)

func (u *UserPostgresql) Insert(user *model.User) error {
	if len(user.ID) < 1 {
		user.Common = model.NewCommon()
	}

	return mdb.DB().WithTx(func(conn mdal.Conn) error {
		_, err := conn.PostgresqlConn.Model(user).Insert()
		return err
	})
}

func (u *UserPostgresql) Query(userIDs []string) ([]*model.User, error) {
	users := make([]*model.User, 0)

	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(&users).Where(model.Common_ID+" in (?)", pg.In(userIDs)).Select()
	})
	if err != nil {
		users = nil
	}

	return users, err
}

func (u *UserPostgresql) QueryOne(userID string) (*model.User, error) {
	user := &model.User{}

	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(user).
			Where(model.User_IsLocked+" = ?", false).
			Where(model.Common_ID+" = ?", userID).
			First()
	})
	if err != nil {
		user = nil
	}

	return user, err
}

func (u *UserPostgresql) QueryOneByUserName(userName string) (*model.User, error) {
	user := &model.User{}

	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(user).
			Where(model.User_IsLocked+" = ?", false).
			Where(model.User_UserName+" = ?", userName).
			First()
	})
	if err != nil {
		user = nil
	}

	return user, err
}

// QueryPageLEPermission
/**
Core: sub-query
	select *
	from users u
	where "permission" <= (
		select "permission"
		from users u
		where id = 'user id'
	) and id != 'user id';
*/
func (u *UserPostgresql) QueryPageLEPermission(
	pageSize int,
	pageNum int,
	userID string,
) ([]*model.User, int, error) {
	users := make([]*model.User, 0)
	count := 0

	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		permission := conn.PostgresqlConn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.Common_ID+" = ?", userID)

		var err2 error
		count, err2 = conn.PostgresqlConn.Model(&users).
			Where(model.User_Permission+" <= (?)", permission).
			Where(model.Common_ID+" != ?", userID).
			Order(model.User_Permission + " DESC").
			Offset((pageNum - 1) * pageSize).
			Limit(pageSize).
			SelectAndCount()

		return err2
	})
	if err != nil {
		users = nil
		count = 0
	}

	return users, count, err
}

func (u *UserPostgresql) UpdateColumnsByUserID(user *model.User, columns ...string) error {
	user.UpdateTime = time.Duration(time.Now().Unix())

	return mdb.DB().WithTx(func(conn mdal.Conn) error {
		query := conn.PostgresqlConn.Model(user).Column(model.Common_UpdateTime)
		for i := range columns {
			query.Column(columns[i])
		}

		_, err := query.Where(model.Common_ID + " = ?" + model.Common_ID).Update()
		return err
	})
}
