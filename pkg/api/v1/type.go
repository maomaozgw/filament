package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type Type struct {
	typeDa *da.FilamentType
	*api.SimpleRest[model.Type]
}

func NewType(factory *da.Factory) (api.RestAPI, error) {
	typeDa, err := da.Get[da.FilamentType](factory)
	if err != nil {
		return nil, err
	}
	api, err := api.NewSimpleRest[model.Type]("/types", api.WithList(typeDa.Search))
	if err != nil {
		return nil, err
	}
	t := &Type{
		typeDa:     typeDa,
		SimpleRest: api,
	}
	return t, nil
}

// List
// @Summary      List Types with filter and page
// @Description  get types by filter
// @Tags         MetaData
// @Accept       json
// @Produce      json
// @Param        page   query      int  true  "Page Number" minimum(1)
// @Param        size   query      int  true  "Page Size" minimum(10)    maximum(100)
// @Param        name   query      string  true  "Filament Type"
// @Param        major query string false "Filament Major Type"
// @Param minor query string false "Filament Minor Type"
// @Success      200  {object}  api.ListResponse[model.Type]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/types [get]
func (t *Type) List(c *gin.Context) {
	t.SimpleRest.List(c)
}

// Create
// @Summary      Create Type
// @Description  Create type
// @Tags         MetaData
// @Accept       json
// @Produce      json
// @Param        type   body      model.Type  true  "Type"
// @Success      200  {object}  api.ItemResponse[model.Type]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/types [post]
func (t *Type) Create(c *gin.Context) {
	t.SimpleRest.Create(c)
}

// Update
// @Summary      Update Type
// @Description  Update type
// @Tags         MetaData
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Brand ID"
// @Param        type   body      model.Type  true  "Type"
// @Success      200  {object}  api.ItemResponse[model.Type]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/types/{id} [put]
func (t *Type) Update(c *gin.Context) {
	t.SimpleRest.Update(c)
}
