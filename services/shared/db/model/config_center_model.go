package model

type Version struct {
	VersionNum  string `pg:",unique,notnull"` // format: v1.2.3-alpha
	Description string

	Services       []string // support service ids
	Configurations []string // service contains config ids

	Common
}
