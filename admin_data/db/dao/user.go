package dao

import (
	"github.com/go-pg/pg/v10/orm"
	. "github.com/mats9693/unnamed_plan/admin_data/db"
	"github.com/mats9693/unnamed_plan/admin_data/db/model"
	"github.com/mats9693/utils/uuid"
)

type User model.User

var userIns = &User{}

func GetUser() *User {
	return userIns
}

func (u *User) Insert(name, password string, permission uint8) (err error) {
	user := &model.User{
		UserID:     uuid.New(),
		Name:       name,
		Password:   password,
		Permission: permission,
	}

	err = WithTx(func(conn orm.DB) error {
		_, err = conn.Model().Insert(user)
		return err
	})

	return err
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
