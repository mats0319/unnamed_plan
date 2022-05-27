package mtest

import (
	"fmt"
	"github.com/go-pg/pg/v10/orm"
	"github.com/mats9693/unnamed_plan/services/shared/db"
	"github.com/mats9693/unnamed_plan/services/shared/db/dal"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"go.uber.org/zap"
	"reflect"
	"time"
)

func init() {
	timestamp := time.Now().Unix()
	// test table have '_[timestamp]' suffix which can distinguish with formal table
	orm.SetTableNameInflector(func(s string) string {
		return fmt.Sprintf("%s_%d", s, timestamp) // make all tables have same timestamp suffix
	})
}

// CreateTestTable_Postgresql pre set test data
//   @param models: e.g. (*model.User)(nil)
func CreateTestTable_Postgresql(models []interface{}) error {
	err := mdb.DB().WithTx(func(conn mdal.Conn) error {
		for i, m := range models {
			err := conn.PostgresqlConn.Model(m).CreateTable(&orm.CreateTableOptions{Temp: false, IfNotExists: true})
			if err != nil {
				mlog.Logger().Error("create table failed",
					zap.Int("index", i),
					zap.String("name", reflect.TypeOf(m).Name()),
					zap.Error(err))
				return err
			}
		}

		return nil
	})
	if err != nil {
		mlog.Logger().Error("create test table failed", zap.Error(err))
		return err
	}

	return nil
}

func DropTestTable_Postgresql(models []interface{}) {
	_ = mdb.DB().WithNoTx(func(conn mdal.Conn) error {
		for i, m := range models {
			err := conn.PostgresqlConn.Model(m).DropTable(&orm.DropTableOptions{
				IfExists: true,
			})
			if err != nil {
				// if a table drop failed, just log it and continue
				mlog.Logger().Error("drop table failed",
					zap.Int("index", i),
					zap.String("name", reflect.TypeOf(m).Name()),
					zap.Error(err))
			}
		}

		return nil
	})
}
