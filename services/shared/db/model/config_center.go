package model

type Administrator struct {
	UserName string `pg:",unique,notnull"`
	Password string `pg:"type:varchar(64),notnull"`
	Salt     string `pg:"type:varchar(10),notnull"`

	Common
}

type Version struct {
	VersionNum  string `pg:",unique,notnull"` // format: v1.2.3-alpha
	Description string

	ServiceIDs     []string // support service ids
	ConfigIDs      []string
	Configurations []string // service contains configs payload(json str)

	IsUsing   bool `pg:",use_zero,notnull"`
	HasUpdate bool `pg:",use_zero,notnull"` // after apply this version, supported service or config has an update

	Common
}

type Service struct {
	ServiceID   string `pg:",unique,notnull"` // generate uuid by service name
	ServiceName string `pg:",unique,notnull"`
	ConfigIDs   []string

	IsShadow bool `pg:",use_zero,notnull"` // shadow service will be skipped when generate new version

	Common
}

type Config struct {
	ConfigID   string `pg:",unique,notnull"` // generate uuid by config name
	ConfigName string `pg:",unique,notnull"`
	Payload    string `pg:",unique,notnull"` // json str

	IsShadow bool `pg:",use_zero,notnull"` // shadow config can not mount to service

	Common
}
