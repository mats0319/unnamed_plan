package mdb

import (
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
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
		ins.postgresql.AddQueryHook(pgdebug.DebugHook{
			Verbose: true,
		})
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
