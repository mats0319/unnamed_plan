package mdal

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mats9693/unnamed_plan/services/shared/log"
	"time"
)

type postgresqlDB struct {
	postgresql *pg.DB
}

var postgresqlDBIns = &postgresqlDB{}

func InitPostgresqlDB(addr string, dbName string, user string, password string, timeout int, showSQL bool) DAL {
	if postgresqlDBIns.postgresql != nil { // have initialized
		return postgresqlDBIns
	}

	postgresqlDBIns.postgresql = pg.Connect(&pg.Options{
		Addr:         addr,
		Database:     dbName,
		User:         user,
		Password:     password,
		ReadTimeout:  time.Duration(timeout) * time.Second, // default behavior is blocking
		WriteTimeout: time.Duration(timeout) * time.Second, // default behavior is blocking
	})

	if showSQL {
		postgresqlDBIns.postgresql.AddQueryHook(&debugHook{})
	}

	return postgresqlDBIns
}

// WithTx
// 其实函数有两个error应该返回，一个是task函数的error，另一个是事务提交或回滚的error
// 此处忽略事务的error、返回task函数的error，主要考虑是task函数出错概率更高，且我们更需要它来debug
func (db *postgresqlDB) WithTx(task func(conn Conn) error) error {
	if task == nil {
		return nil
	}

	conn := db.postgresqlConn()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	err = task(Conn{PostgresqlConn: tx})

	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}

	return err
}

func (db *postgresqlDB) WithNoTx(task func(conn Conn) error) error {
	if task == nil {
		return nil
	}

	conn := db.postgresqlConn()
	defer conn.Close()

	return task(Conn{PostgresqlConn: conn})
}

func (db *postgresqlDB) postgresqlConn() *pg.Conn {
	return db.postgresql.Conn()
}

// debugHook simulate pgdebug.DebugHook, for use self-define print method, redirect SQL into log file
type debugHook struct {
}

var _ pg.QueryHook = (*debugHook)(nil)

// BeforeQuery only record SQL, ignore errors this duration
func (d *debugHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	queryBytes, err := event.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if event.Err != nil {
		mlog.Logger().Info(fmt.Sprintf("%s executing a query:\n%s\n", event.Err, queryBytes))
	} else {
		mlog.Logger().Info(string(queryBytes))
	}

	return ctx, nil
}

func (d *debugHook) AfterQuery(_ context.Context, _ *pg.QueryEvent) error {
	return nil
}
