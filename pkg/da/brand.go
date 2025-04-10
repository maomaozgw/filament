package da

import (
	"context"
	"errors"

	"github.com/maomaozgw/filament/pkg/model"
	"gorm.io/gorm"
)

func init() {
	typeFactory[NameOf[Brand]()] = func(db *gorm.DB) any {
		return NewBrand(db)
	}
}

type Brand struct {
	db *gorm.DB
}

func NewBrand(db *gorm.DB) *Brand {
	return &Brand{
		db: db,
	}
}

func (b *Brand) ListAll(ctx context.Context) ([]model.Brand, error) {
	var result []model.Brand
	err := b.db.WithContext(ctx).Find(&result).Error
	return result, err
}

func (b *Brand) Search(ctx context.Context, filter map[string]string, page, pageSize int) ([]model.Brand, int64, error) {
	db := b.db.WithContext(ctx).Model(model.Brand{})
	var result []model.Brand
	var total int64
	var err error
	if name, ok := filter["name"]; ok {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if db = db.Count(&total); db.Error != nil {
		return nil, 0, db.Error
	}

	err = db.Limit(pageSize).Offset(pageSize * (page - 1)).Find(&result).Error
	return result, total, err
}

func (b *Brand) get(ctx context.Context, db *gorm.DB, id uint, name string) (*model.Brand, error) {
	var brand = &model.Brand{}
	db = db.WithContext(ctx)
	if id > 0 {
		db = db.First(brand, id)
	} else {
		db = db.Model(brand).Where("name = ?", name).First(brand)
	}
	return brand, db.Error
}

func (b *Brand) getOrCreate(ctx context.Context, db *gorm.DB, name string) (*model.Brand, error) {
	var err error
	var info = &model.Brand{
		Name: name,
	}
	result, err := b.get(ctx, db, 0, name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return info, b.create(ctx, db, info)
	} else if err != nil {
		return nil, err
	}
	return result, err
}

func (b *Brand) create(ctx context.Context, db *gorm.DB, brand *model.Brand) error {
	return db.WithContext(ctx).Create(brand).Error
}

func (b *Brand) Get(ctx context.Context, id uint) (*model.Brand, error) {
	return b.get(ctx, b.db, id, "")
}

func (b *Brand) GetByName(ctx context.Context, name string) (*model.Brand, error) {
	return b.get(ctx, b.db, 0, name)
}

func (b *Brand) Create(ctx context.Context, brand *model.Brand) (*model.Brand, error) {
	return b.getOrCreate(ctx, b.db, brand.Name)
}

func (b *Brand) Update(ctx context.Context, id uint, brand *model.Brand) (*model.Brand, error) {
	var result *model.Brand
	err := b.db.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = b.get(ctx, tx, id, "")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Name = brand.Name
			return b.create(ctx, tx, result)
		} else if err != nil {
			return err
		}
		result.Name = brand.Name
		return tx.Save(result).Error
	})
	return result, err
}
