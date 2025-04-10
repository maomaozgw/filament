package server

type Option struct {
	Addr string `mapstructure:"addr"`
	Orm  OrmOpt `mapstructure:"store"`
}

const (
	StoreTypeSqlite = "sqlite"
)

type OrmOpt struct {
	Type   string     `mapstructure:"type"`
	Sqlite *SqliteOpt `mapstructure:"sqlite"`
}

type SqliteOpt struct {
	Path string `mapstructure:"path"`
}
