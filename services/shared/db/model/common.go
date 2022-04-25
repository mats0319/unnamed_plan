package model

import (
	"github.com/mats9693/utils/uuid"
	"time"
)

type Common struct {
	ID             string        `pg:",pk"`
	CreatedTime    time.Duration `pg:",use_zero,notnull"`
	UpdateTime     time.Duration `pg:",use_zero,notnull"`
	OptimisticLock uint64        `pg:",use_zero,notnull"`
}

func NewCommon() Common {
	currTime := time.Duration(time.Now().Unix())

	return Common{
		ID:             uuid.New(),
		UpdateTime:     currTime,
		CreatedTime:    currTime,
		OptimisticLock: 0,
	}
}
