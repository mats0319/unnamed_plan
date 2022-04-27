package dao

import (
    "github.com/mats9693/unnamed_plan/services/shared/db"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/db/model"
)

type ServiceConfigPostgresql model.ServiceConfig

var _ ServiceConfigDao = (*ServiceConfigPostgresql)(nil)

func (s ServiceConfigPostgresql) Query() ([]*model.ServiceConfig, error) {
    config := make([]*model.ServiceConfig, 0)

    err := mdb.DB().WithNoTx(func(conn mdal.Conn) error {
        return conn.PostgresqlConn.Model(&config).Where(model.ServiceConfig_IsDelete+" = ?", false).Select()
    })
    if err != nil {
        config = nil
    }

    return config, err
}
