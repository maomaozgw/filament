package da

import (
	"context"
	"log"
	"sort"

	"github.com/maomaozgw/filament/pkg/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func init() {
	typeFactory[NameOf[Warehouse]()] = func(db *gorm.DB) any {
		return NewInventory(db)
	}
}

type Warehouse struct {
	db           *gorm.DB
	brand        *Brand
	filamentType *FilamentType
	color        *Color
}

func NewInventory(db *gorm.DB) *Warehouse {
	return &Warehouse{
		db:           db,
		brand:        NewBrand(db),
		filamentType: NewFilamentType(db),
		color:        NewColor(db),
	}
}

func (i *Warehouse) Import(ctx context.Context, items []model.Filament) error {
	db := i.db.WithContext(ctx)
	// prepare data
	var brandMap = map[string]uint{}
	var colorMap = map[string]uint{}
	var typeMap = map[string]uint{}
	for _, item := range items {
		item.ID = 0
		brandMap[item.Brand.Name] = 0
		colorMap[item.Color.Name] = 0
		typeMap[item.Type.Name] = 0
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for name := range brandMap {
			brand, err := i.brand.getOrCreate(ctx, tx, name)
			if err != nil {
				return err
			}
			brandMap[name] = brand.ID
		}
		for name := range colorMap {
			color, err := i.color.getOrCreate(ctx, tx, &model.Color{Name: name})
			if err != nil {
				return err
			}
			colorMap[name] = color.ID
		}
		for name := range typeMap {
			typ, err := i.filamentType.getOrCreate(ctx, tx, &model.Type{Name: name})
			if err != nil {
				return err
			}
			typeMap[name] = typ.ID
		}
		for idx := range items {
			item := items[idx]
			item.ID = 0
			item.BrandId = brandMap[item.Brand.Name]
			item.Brand.ID = brandMap[item.Brand.Name]
			item.ColorId = colorMap[item.Color.Name]
			item.Color.ID = colorMap[item.Color.Name]
			item.TypeId = typeMap[item.Type.Name]
			item.Type.ID = typeMap[item.Type.Name]
			if err := i.upsert(ctx, tx, &item); err != nil {
				return err
			}
		}
		return nil
	})
}

func (i *Warehouse) upsert(ctx context.Context, tx *gorm.DB, info *model.Filament) error {
	inv, err := i.get(ctx, tx, info)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("create filament %+v", info)
		err = tx.Create(&info).Error
	} else if err != nil {
		return err
	} else {
		inv.Quantity += info.Quantity
		err = tx.Save(&inv).Error
	}
	if err != nil {
		return err
	}
	err = tx.Create(&model.Record{
		BrandId:  info.BrandId,
		TypeId:   info.TypeId,
		ColorId:  info.ColorId,
		Quantity: info.Quantity,
		Price:    info.Price,
		Kind:     model.KindStockIn,
	}).Error
	return err
}

func (i *Warehouse) StockIn(ctx context.Context, info *model.Filament) error {
	db := i.db.WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		if info.BrandId == 0 {
			brand, err := i.brand.getOrCreate(ctx, tx, info.Brand.Name)
			if err != nil {
				return err
			}
			info.BrandId = brand.ID
		}
		if info.TypeId == 0 {
			t, err := i.filamentType.getOrCreate(ctx, tx, &info.Type)
			if err != nil {
				return err
			}
			info.TypeId = t.ID
		}
		if info.ColorId == 0 {
			c, err := i.color.getOrCreate(ctx, tx, &info.Color)
			if err != nil {
				return err
			}
			info.ColorId = c.ID
		}
		return i.upsert(ctx, tx, info)
	})
}

func (i *Warehouse) get(ctx context.Context, db *gorm.DB, info *model.Filament) (*model.Filament, error) {
	var result = &model.Filament{}
	db = db.WithContext(ctx).Preload("Brand").Preload("Type").Preload("Color")
	if info.ID > 0 {
		v := db.Where("id = ?", info.ID).First(&result)
		return result, v.Error
	}
	if info.BrandId > 0 {
		db = db.Where("brand_id =?", info.BrandId)
	} else {
		db = db.Joins("JOIN brands on brands.id = brand_id").Where("brands.name =?", info.Brand.Name)
	}
	if info.TypeId > 0 {
		db = db.Where("type_id =?", info.TypeId)
	} else {
		major, minor := model.ExploreType(info.Type.Name)
		db = db.Joins("JOIN types on types.id = type_id").Where("types.major =? and types.minor = ?", major, minor)
	}
	if info.ColorId > 0 {
		db = db.Where("color_id =?", info.ColorId)
	} else {
		db = db.Joins("JOIN colors on colors.id = color_id").Where("colors.name =?", info.Color.Name)
	}
	v := db.First(&result)
	return result, v.Error
}

func (i *Warehouse) StockTake(ctx context.Context, info *model.Filament) error {
	db := i.db.WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		inv, err := i.get(ctx, tx, info)
		if err != nil {
			return err
		}
		inv.Quantity = info.Quantity
		if err = tx.Save(&inv).Error; err != nil {
			return err
		}
		err = tx.Create(&model.Record{
			BrandId:  info.BrandId,
			TypeId:   info.TypeId,
			Quantity: inv.Quantity - info.Quantity,
			Kind:     model.KindStockTake,
		}).Error
		return err
	})
}

func (i *Warehouse) StockOut(ctx context.Context, info *model.Filament) error {
	db := i.db.WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		inv, err := i.get(ctx, tx, info)
		if err != nil {
			return err
		}
		inv.Quantity -= info.Quantity
		if inv.Quantity < 0 {
			return errors.New("not enough filament")
		}
		err = tx.Save(&inv).Error
		if err != nil {
			return err
		}
		err = tx.Create(&model.Record{
			BrandId:  info.BrandId,
			TypeId:   info.TypeId,
			Quantity: info.Quantity,
			Kind:     model.KindStockOut,
		}).Error
		return err
	})
}

func (i *Warehouse) SearchStatistic(ctx context.Context, filter map[string]string, page, pageSize int) ([]model.Statistic, int64, error) {
	db := i.db.WithContext(ctx)
	var result []model.Statistic
	summary := model.Statistic{
		Kind:  "summary",
		Title: "Summary",
	}
	totalCount := model.StatisticValue{Name: "Current"}
	if err := db.Model(model.Filament{}).Select("sum(quantity) as value").First(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	totalStockIn := model.StatisticValue{Name: "Stock In"}
	if err := db.Model(model.Record{}).Where("kind =?", model.KindStockIn).Select("sum(quantity) as value").First(&totalStockIn).Error; err != nil {
		return nil, 0, err
	}
	totalStockOut := model.StatisticValue{Name: "Stock Out"}
	if err := db.Model(model.Record{}).Where("kind <> ?", model.KindStockIn).Select("sum(quantity) as value").First(&totalStockOut).Error; err != nil {
		return nil, 0, err
	}
	totalCost := model.StatisticValue{Name: "Total Cost"}
	if err := db.Model(model.Record{}).Select("sum(quantity * price) as value").First(&totalCost).Error; err != nil {
		return nil, 0, err
	}
	summary.Values = append(summary.Values, totalCount, totalStockIn, totalStockOut, totalCost)

	brandAgg := []model.StatisticValue{}
	if err := db.Model(model.Filament{}).Joins("JOIN brands on brands.id = brand_id").Select("brands.name as name, sum(quantity) as value").Group("brand_id").Find(&brandAgg).Error; err != nil {
		return nil, 0, err
	}
	typeTmpAgg := []struct {
		Major string `json:"major"`
		Minor string `json:"minor"`
		Value int64  `json:"value"`
	}{}
	if err := db.Model(model.Filament{}).Joins("JOIN types on types.id = type_id").Select("types.major,types.minor, sum(quantity) as value").Group("type_id").Find(&typeTmpAgg).Error; err != nil {
		return nil, 0, err
	}
	typeAggMap := map[string]*model.StatisticValue{}
	for _, v := range typeTmpAgg {
		major, minor := v.Major, v.Minor
		val := v.Value
		if _, ok := typeAggMap[major]; !ok {
			typeAggMap[major] = &model.StatisticValue{
				Name:     major,
				Chileren: []model.StatisticValue{},
			}
		}
		typeAggMap[major].Chileren = append(typeAggMap[major].Chileren, model.StatisticValue{
			Name:  minor,
			Value: val,
		})
		typeAggMap[major].Value += val
	}
	typeAgg := []model.StatisticValue{}
	for _, v := range typeAggMap {
		typeAgg = append(typeAgg, *v)
	}
	sort.Slice(typeAgg, func(i, j int) bool {
		return typeAgg[i].Value > typeAgg[j].Value
	})
	colorAgg := []model.StatisticValue{}
	if err := db.Model(model.Filament{}).Joins("JOIN colors on colors.id = color_id").Select("colors.name as name, sum(quantity) as value").Group("color_id").Find(&colorAgg).Error; err != nil {
		return nil, 0, err
	}
	result = append(result, summary, model.Statistic{
		Kind:   "Pie",
		Title:  "Brand",
		Values: brandAgg,
	}, model.Statistic{
		Kind:   "Sunburst",
		Title:  "Type",
		Values: typeAgg,
	}, model.Statistic{
		Kind:   "Pie",
		Title:  "Color",
		Values: colorAgg,
	})
	return result, 0, nil
}

func (i *Warehouse) SearchWarehouse(ctx context.Context, filter map[string]string, page, pageSize int) ([]model.Filament, int64, error) {
	db := i.db.WithContext(ctx).Model(model.Filament{})
	var result []model.Filament
	var total int64
	var err error
	db, err = i.buildQuery(ctx, db, filter)
	if err != nil {
		return nil, 0, err
	}
	db = db.Count(&total)
	if db.Error != nil {
		return nil, 0, db.Error
	}
	err = db.Preload("Brand").Preload("Type").Preload("Color").Limit(pageSize).Offset(pageSize * (page - 1)).Find(&result).Error
	return result, total, err
}

func (i *Warehouse) SearchRecord(ctx context.Context, filter map[string]string, page, pageSize int) ([]model.Record, int64, error) {
	db := i.db.WithContext(ctx).Model(model.Record{})
	var result []model.Record
	var total int64
	var err error
	db, err = i.buildQuery(ctx, db, filter)
	if err != nil {
		return nil, 0, err
	}
	if db = db.Count(&total); db.Error != nil {
		return nil, 0, db.Error
	}
	err = db.Preload("Brand").Preload("Type").Preload("Color").Limit(pageSize).Offset(pageSize * (page - 1)).Find(&result).Error
	return result, total, err

}

func (i *Warehouse) GetFilament(ctx context.Context, id uint) (*model.Filament, error) {
	return i.get(
		ctx, i.db,
		&model.Filament{
			Base: model.Base{ID: id},
		},
	)
}

func (i *Warehouse) buildQuery(ctx context.Context, db *gorm.DB, filter map[string]string) (*gorm.DB, error) {
	const (
		Type      = "type"
		BrandId   = "brand_id"
		Brand     = "brand"
		TypeId    = "type_id"
		ColorId   = "color_id"
		Color     = "color"
		TypeMajor = "type_major"
		TypeMinor = "type_minor"
		Kind      = "kind"
	)
	if k, ok := filter[Kind]; ok {
		db = db.Where("kind = ?", k)
	}
	if t, ok := filter["type"]; ok {
		filter[TypeMajor], filter[TypeMinor] = model.ExploreType(t)
	}
	if id, ok := filter[BrandId]; ok {
		db = db.Where("brand_id =?", id)
	} else if brand, ok := filter[Brand]; ok {
		b, err := i.brand.get(ctx, db, 0, brand)
		if err != nil {
			return nil, errors.Wrapf(err, "get brand %s failed", brand)
		}
		db = db.Where("brand_id =?", b.ID)
	}
	if id, ok := filter[TypeId]; ok {
		db = db.Where("type_id =?", id)
	} else if major, ok := filter[TypeMajor]; ok {
		db.Joins("JOIN types on types.id = type_id").Where("types.major =?", major)
	}
	if id, ok := filter[ColorId]; ok {
		db = db.Where("color_id =?", id)
	} else if color, ok := filter[Color]; ok {
		c, err := i.color.get(ctx, db, 0, color)
		if err != nil {
			return nil, errors.Wrapf(err, "get color %s failed", color)
		}
		db = db.Where("color_id =?", c.ID)
	}
	return db, nil
}
