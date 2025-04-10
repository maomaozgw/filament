package da

import (
	"context"
	"errors"

	"github.com/maomaozgw/filament/pkg/model"
	"gorm.io/gorm"
)

func init() {
	typeFactory[NameOf[FilamentType]()] = func(db *gorm.DB) any {
		return NewFilamentType(db)
	}
}

type Filter struct {
	Brand     string
	Type      string
	TypeMajor string
	TypeMinor string
	Color     string

	BrandId uint
	TypeId  uint
	ColorId uint
}

type FilamentType struct {
	db *gorm.DB
}

func NewFilamentType(db *gorm.DB) *FilamentType {
	return &FilamentType{db: db}
}

func (f *FilamentType) ListAll(ctx context.Context) ([]model.Type, error) {
	var result []model.Type
	err := f.db.WithContext(ctx).Find(&result).Error
	return result, err
}

func (f *FilamentType) getOrCreate(ctx context.Context, db *gorm.DB, info *model.Type) (*model.Type, error) {
	result, err := f.get(ctx, db, info.ID, info.Name, info.Major, info.Minor)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return info, f.create(ctx, db, info)
	} else if err != nil {
		return nil, err
	}
	return result, err
}

func (f *FilamentType) create(ctx context.Context, db *gorm.DB, info *model.Type) error {
	return db.WithContext(ctx).Create(info).Error
}

func (f *FilamentType) get(ctx context.Context, db *gorm.DB, id uint, name, major, minor string) (*model.Type, error) {
	var result = &model.Type{}
	db = db.WithContext(ctx)
	if id > 0 {
		v := db.First(&result)
		return result, v.Error
	} else if name != "" {
		major, minor := model.ExploreType(name)
		v := db.Where("major =?", major).Where("minor =?", minor).First(&result)
		return result, v.Error
	} else if major != "" && minor != "" {
		v := db.Where("major =?", major).Where("minor =?", minor).First(&result)
		return result, v.Error
	} else {
		return nil, errors.New("invalid get request")
	}

}

func (f *FilamentType) Search(ctx context.Context, filter map[string]string, page, pageSize int) ([]model.Type, int64, error) {
	var result []model.Type
	var total int64
	var err error
	db := f.db.WithContext(ctx).Model(model.Type{})
	if major, ok := filter["major"]; ok {
		db = db.Where("major =?", major)
	}
	if minor, ok := filter["minor"]; ok {
		db = db.Where("minor =?", minor)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err = db.Limit(pageSize).Offset(pageSize * (page - 1)).Find(&result).Error
	return result, total, err
}
