package dao

import (
	"github.com/go-pg/pg/v10/orm"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	. "github.com/mats9693/unnamed_plan/shared/go/db"
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
		return conn.Model(user).Where("is_locked = ?", false).Where(condition, param...).First()
	})
	if err != nil {
		user = nil
	}

	return
}

func (u *User) QueryPage(
	pageSize int,
	pageNum int,
	condition string,
	param ...interface{},
) (users []*model.User, count int, err error) {
	err = WithNoTx(func(conn orm.DB) error {
		count, err = conn.Model(&users).Where(condition, param...).Order("permission DESC").
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
		query := conn.Model(data).Column("update_time")
		for i := range columns {
			query.Column(columns[i])
		}

		_, err = query.Where("user_id = ?user_id").Update()
		return err
	})
}
