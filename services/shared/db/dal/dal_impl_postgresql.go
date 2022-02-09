package mdb

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

var _ DAL = (*postgresqlDB)(nil)

func initPostgresqlDB(addr string, dbName string, user string, password string, timeout int, showSQL bool) DAL {
	ins := &postgresqlDB{}

	ins.postgresql = pg.Connect(&pg.Options{
		Addr:         addr,
		Database:     dbName,
		User:         user,
		Password:     password,
		ReadTimeout:  time.Duration(timeout) * time.Second, // default behavior is blocking
		WriteTimeout: time.Duration(timeout) * time.Second, // default behavior is blocking
	})

	if showSQL {
		ins.postgresql.AddQueryHook(&debugHook{})
	}

	return ins
}

func (db *postgresqlDB) WithTx(task func(conn Conn) error) error {
	if task == nil {
		return nil
	}

	conn := db.postgresqlConn()
	defer func() {
		_ = conn.Close()
	}()

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			err = tx.Rollback()
		}
	}()

	err = task(Conn{PostgresqlConn: tx})

	return err
}

func (db *postgresqlDB) WithNoTx(task func(conn Conn) error) error {
	if task == nil {
		return nil
	}

	conn := db.postgresqlConn()
	defer func() {
		_ = conn.Close()
	}()

	err := task(Conn{PostgresqlConn: conn})

	return err
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
