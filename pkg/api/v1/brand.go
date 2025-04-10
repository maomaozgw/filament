package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type Brand struct {
	brand *da.Brand
	*api.SimpleRest[model.Brand]
}

func NewBrand(factory *da.Factory) (api.RestAPI, error) {
	brandDa, err := da.Get[da.Brand](factory)
	if err != nil {
		return nil, err
	}
	simpleRest, err := api.NewSimpleRest[model.Brand](
		"/brands",
		api.WithList(brandDa.Search),
		api.WithCreate(brandDa.Create),
		api.WithUpdate(brandDa.Update),
	)
	b := &Brand{
		brand:      brandDa,
		SimpleRest: simpleRest,
	}
	return b, nil
}

// List
// @Summary List Brand
// @Description  List brands with filter and pagenation
// @Tags MetaData
// @Accept       json
// @Produce      json
// @Param        page   query      int  true  "Page Number" minimum(1)
// @Param        size   query      int  true  "Page Size" minimum(10)    maximum(100)
// @Param        name   query      string  true  "Brand Name"
// @Success      200  {object}  api.ListResponse[model.Brand]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/brands [get]
func (b *Brand) List(c *gin.Context) {
	b.SimpleRest.List(c)
}

// Create
// @Summary Create Brand
// @Description  Create brand
// @Tags MetaData
// @Accept       json
// @Produce      json
// @Param        brand   body      model.Brand  true  "Brand"
// @Success      200  {object}  api.ItemResponse[model.Brand]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/brands [post]
func (b *Brand) Create(c *gin.Context) {
	b.SimpleRest.Create(c)
}

// Update
// @Summary Update Brand
// @Description  Update brand
// @Tags MetaData
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Brand ID"
// @Param        brand   body      model.Brand  true  "Brand"
// @Success      200  {object}  api.ItemResponse[model.Brand]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/brands/{id} [put]
func (b *Brand) Update(c *gin.Context) {
	b.SimpleRest.Update(c)
}
