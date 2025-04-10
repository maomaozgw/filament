package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type metaData struct {
	Colors []model.Color `json:"colors"`
	Brands []model.Brand `json:"brands"`
	Types  []model.Type  `json:"types"`
}

type MetaData struct {
	Color *da.Color
	Brand *da.Brand
	Type  *da.FilamentType
}

// RegisterRoutes implements api.RestAPI.
func (m *MetaData) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("", api.RestWrapper(m.Get))
}

func NewMetaData(factory *da.Factory) (api.RestAPI, error) {
	color, err := da.Get[da.Color](factory)
	if err != nil {
		return nil, err
	}
	brand, err := da.Get[da.Brand](factory)
	if err != nil {
		return nil, err
	}
	filamentType, err := da.Get[da.FilamentType](factory)
	if err != nil {
		return nil, err
	}
	m := &MetaData{
		Color: color,
		Brand: brand,
		Type:  filamentType,
	}
	return m, nil
}

// Get
// @Summary Get all metadata
// @Description Get Global Metadata, including colors, brands, and types.
// @Tags MetaData
// @Accept json
// @Produce json
// @Success 200 {object} api.ItemResponse[metaData]
// @Router /v1/meta-data [get]
func (m *MetaData) Get(c *gin.Context) (*api.ItemResponse[metaData], error) {
	colors, _, err := m.Color.List(c, nil, 1, 10000)
	if err != nil {
		return nil, err
	}
	brands, err := m.Brand.ListAll(c)
	if err != nil {
		return nil, err
	}
	types, err := m.Type.ListAll(c)
	if err != nil {
		return nil, err
	}
	return &api.ItemResponse[metaData]{
		Data: &metaData{
			Colors: colors,
			Brands: brands,
			Types:  types,
		},
	}, nil
}
