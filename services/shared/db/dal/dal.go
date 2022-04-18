package mdal

import "github.com/go-pg/pg/v10/orm"

// DAL data access layer
type DAL interface {
	WithTx(func(Conn) error) error
	WithNoTx(func(Conn) error) error
}

type Conn struct {
	PostgresqlConn orm.DB
}
