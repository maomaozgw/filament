package server

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newOrm(opt OrmOpt) (*gorm.DB, error) {
	switch opt.Type {
	case StoreTypeSqlite:
		return gorm.Open(sqlite.Open(opt.Sqlite.Path), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})

	default:
		return nil, nil
	}
}
