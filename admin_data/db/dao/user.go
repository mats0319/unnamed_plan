package dao

import (
	"github.com/go-pg/pg/v10/orm"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	. "github.com/mats9693/utils/toy_server/db"
	"github.com/mats9693/utils/uuid"
	"time"
)

type User model.User

var userIns = &User{}

func GetUser() *User {
	return userIns
}

func (u *User) Insert(user *model.User) (err error) {
	if len(user.ID) < 1 {
		user.Common = model.NewCommon()
	}

	if len(user.UserID) < 1 {
		user.UserID = uuid.New()
	}

	return WithTx(func(conn orm.DB) error {
		_, err = conn.Model(user).Insert()
		return err
	})
}

// QueryOne only query users unlocked
func (u *User) QueryOne(condition string, param ...interface{}) (user *model.User, err error) {
	user = &model.User{}

	err = WithNoTx(func(conn orm.DB) error {
		return conn.Model(user).Where(model.User_IsLocked+" = ?", false).Where(condition, param...).First()
	})
	if err != nil {
		user = nil
	}

	return
}

// QueryPageByPermission 获取用户列表，要求目标用户权限等级不高于指定用户（通过userID指定），分页，按照权限等级降序
/**
Core: sub-query
	select *
	from users u
	where "permission" <= (
		select "permission"
		from users u
		where user_id = 'user id'
	);
*/
func (u *User) QueryPageByPermission(
	pageSize int,
	pageNum int,
	userID string,
) (users []*model.User, count int, err error) {
	err = WithNoTx(func(conn orm.DB) error {
		permission := conn.Model((*model.User)(nil)).Column(model.User_Permission).Where(model.User_UserID+" = ?", userID)

		count, err = conn.Model(&users).Where(model.User_Permission+" <= ?", permission).
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

func (u *User) Query(condition string, param ...interface{}) (users []*model.User, err error) {
	err = WithNoTx(func(conn orm.DB) error {
		return conn.Model(&users).Where(condition, param...).Select()
	})
	if err != nil {
		users = nil
	}

	return
}

func (u *User) UpdateColumnsByUserID(data *model.User, columns ...string) (err error) {
	data.UpdateTime = time.Duration(time.Now().Unix())

	return WithTx(func(conn orm.DB) error {
		query := conn.Model(data).Column(model.Common_UpdateTime)
		for i := range columns {
			query.Column(columns[i])
		}

		_, err = query.Where(model.User_UserID + " = ?" + model.User_UserID).Update()
		return err
	})
}
