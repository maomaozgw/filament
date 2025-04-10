package da

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

var (
	typeFactory = map[string]func(db *gorm.DB) any{}
)

type Factory struct {
	db            *gorm.DB
	createFuncMap map[string]func(db *gorm.DB) any
}

func NewFactory(db *gorm.DB) (*Factory, error) {
	return &Factory{
		db:            db,
		createFuncMap: typeFactory,
	}, nil
}

func NameOf[T any]() string {
	var t T
	return reflect.TypeOf(t).Name()
}

func Register[T any](factory *Factory, createFunc func(db *gorm.DB) any) {
	factory.createFuncMap[NameOf[T]()] = createFunc
}

func MustGet[T any](factory *Factory) *T {
	createFunc, ok := factory.createFuncMap[NameOf[T]()]
	if !ok {
		panic("not found")
	}
	return createFunc(factory.db).(*T)
}

func Get[T any](factory *Factory) (*T, error) {
	createFunc, ok := factory.createFuncMap[NameOf[T]()]
	if !ok {
		return nil, fmt.Errorf("lookup create func for %s failed: not found", NameOf[T]())
	}
	return createFunc(factory.db).(*T), nil

}
