package mconst

type HTTPFlags = uint8

const (
	HTTPFlags_MultiLogin_SkipLimit HTTPFlags = iota
	HTTPFlags_MultiLogin_ReSetParams
)

type ConfigLevel = string

const (
	ConfigLevel_Default    ConfigLevel = "default"
	ConfigLevel_Dev                    = "dev"
	ConfigLevel_Production             = "production"
	ConfigLevel_Test                   = "test"
)
