package server

type Option struct {
	Addr      string `mapstructure:"addr"`
	Orm       OrmOpt `mapstructure:"store"`
	StaticDir string `mapstructure:"static-dir"`
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
