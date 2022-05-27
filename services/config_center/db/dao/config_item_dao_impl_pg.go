package dao

import (
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/db/model"
)

type ConfigItemPostgresql model.ConfigItem

var _ ConfigItemDao = (*ConfigItemPostgresql)(nil)

func (c ConfigItemPostgresql) Query() ([]*model.ConfigItem, error) {
	config := make([]*model.ConfigItem, 0)

	err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		return conn.PostgresqlConn.Model(&config).Select()
	})
	if err != nil {
		config = nil
	}

	return config, err
}
