package mtest

import (
    "fmt"
    "github.com/go-pg/pg/v10/orm"
    "github.com/mats9693/unnamed_plan/services/shared/db/dal"
    "github.com/mats9693/unnamed_plan/services/shared/log"
    "go.uber.org/zap"
    "os"
    "reflect"
    "time"
)

func init() {
    orm.SetTableNameInflector(func(s string) string {
        return fmt.Sprintf("%s_%d", s, time.Now().Unix())
    })
}

// CreateTestTable pre set test data
//   @param models: e.g. (*model.User)(nil)
//   @param preSetDataFunc: insert into db
func CreateTestTable(models []interface{}) {
    err := mdb.DB().WithTx(func(conn mdb.Conn) error {
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
        os.Exit(-1)
    }
}

func DropTestTable(models []interface{}) {
    _ = mdb.DB().WithNoTx(func(conn mdb.Conn) error {
        for i, m := range models {
            err := conn.PostgresqlConn.Model(m).DropTable(&orm.DropTableOptions{
                IfExists: true,
            })
            if err != nil {
                mlog.Logger().Error("drop table failed",
                    zap.Int("index", i),
                    zap.String("name", reflect.TypeOf(m).Name()),
                    zap.Error(err))
            }
        }

        return nil
    })
}
