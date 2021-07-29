package dao

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/db/model"
)

func GetUser(db *pg.DB, name, pwd string) (*model.User, error) {
	var (
		users = make([]model.User, 0)
		err   error
	)

	err = db.Model(&users).Where("is_locked = ?", false).Where("name = ?", name).Select()
	if err != nil {
		// todo: write log
		return nil, err
	}

	var user *model.User

	for i := range users {
		if users[i].Password == pwd {
			user = &users[i]
			break
		}
	}

	if user == nil {
		err = errors.New("invalid user name or password")
		return nil, err
	}

	return user, nil
}
