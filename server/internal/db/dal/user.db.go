package dal

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"gorm.io/gorm"
)

func CreateUser(user *model.User) (e *utils.Error) {
	err := User.WithContext(context.TODO()).Create(user)
	if err != nil {
		pge, ok := errors.AsType[*pgconn.PgError](err)
		if ok && pge.Code == "23505" {
			e = utils.ErrUserExist().WithCause(err)
		} else {
			e = utils.ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return
	}

	return
}

func UpdateUser(user *model.User) *utils.Error {
	err := User.WithContext(context.TODO()).Where(User.ID.Eq(user.ID)).Save(user)
	if err != nil {
		e := utils.ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func GetUser(userName string) (record *model.User, e *utils.Error) {
	record, err := User.WithContext(context.TODO()).Where(User.UserName.Eq(userName)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e = utils.ErrUserNotFound().WithCause(err)
		} else {
			e = utils.ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return
	}

	return
}

func ListUsers(pageSize int, pageNum int) (count int64, records []*model.User, e *utils.Error) {
	sql := User.WithContext(context.TODO())

	count, err := sql.Count()
	if err != nil {
		e = utils.ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	records, err = sql.Order(User.UpdatedAt.Desc()).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find()
	if err != nil {
		e = utils.ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
