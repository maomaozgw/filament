package da

import (
	"context"
	"errors"

	"github.com/maomaozgw/filament/pkg/model"
	"gorm.io/gorm"
)

func init() {
	typeFactory[NameOf[Color]()] = func(db *gorm.DB) any {
		return NewColor(db)
	}
}

type Color struct {
	db *gorm.DB
}

func NewColor(db *gorm.DB) *Color {
	return &Color{db: db}
}

func (c *Color) getOrCreate(ctx context.Context, db *gorm.DB, color *model.Color) (*model.Color, error) {
	result, err := c.get(ctx, db, color.ID, color.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return color, c.create(ctx, db, color)
	} else if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Color) create(ctx context.Context, db *gorm.DB, color *model.Color) error {
	return db.WithContext(ctx).Create(color).Error
}

func (c *Color) get(ctx context.Context, db *gorm.DB, id uint, name string) (*model.Color, error) {
	db = db.WithContext(ctx)
	if id > 0 {
		db = db.Where("id = ?", id)
	} else {
		db = db.Where("name =?", name)
	}
	var result = &model.Color{}
	err := db.First(result).Error
	return result, err
}

func (c *Color) List(ctx context.Context, filter map[string]string, page, pageSize int) ([]model.Color, int64, error) {
	db := c.db.WithContext(ctx)
	var result []model.Color
	err := db.Find(&result).Error
	return result, int64(len(result)), err
}

func (c *Color) Create(ctx context.Context, color *model.Color) (*model.Color, error) {
	var err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return c.create(ctx, tx, color)
	})
	return color, err
}

func (c *Color) Update(ctx context.Context, id uint, color *model.Color) (*model.Color, error) {
	var err = c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		oldColor, err := c.get(ctx, tx, id, "")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.create(ctx, tx, color)
		} else if err != nil {
			return err
		}
		oldColor.Name = color.Name
		oldColor.RGBA = color.RGBA
		err = tx.Save(oldColor).Error
		color = oldColor
		return err
	})
	return color, err
}
