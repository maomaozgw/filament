package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type Color struct {
	colorDa *da.Color
	*api.SimpleRest[model.Color]
}

func NewColor(factory *da.Factory) (api.RestAPI, error) {
	color, err := da.Get[da.Color](factory)
	if err != nil {
		return nil, err
	}
	restAPi, err := api.NewSimpleRest(
		"/colors",
		api.WithList(color.List),
		api.WithCreate(color.Create),
		api.WithUpdate(color.Update),
	)
	if err != nil {
		return nil, err
	}
	c := &Color{colorDa: color, SimpleRest: restAPi}
	return c, nil
}

// List
// @Summary      List All Colors
// @Description  get all colors
// @Tags         MetaData
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ListResponse[model.Color]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/colors [get]
func (c *Color) List(ctx *gin.Context) {
	c.SimpleRest.List(ctx)
}

// Create
// @Summary      Create Color
// @Description  Create Color
// @Tags         MetaData
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ItemResponse[model.Color]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/colors [post]
func (c *Color) Create(ctx *gin.Context) {
	c.SimpleRest.Create(ctx)
}

// Update
// @Summary      Update Color
// @Description  Update Color only allow to update name or rgba setting
// @Tags         MetaData
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.ItemResponse[model.Color]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/meta-data/colors/{id} [put]
func (c *Color) Update(ctx *gin.Context) {
	c.SimpleRest.Update(ctx)
}
