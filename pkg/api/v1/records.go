package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type Record struct {
	warehouse *da.Warehouse
	*api.SimpleRest[model.Record]
}

func NewRecords(factory *da.Factory) (api.RestAPI, error) {
	warehouse, err := da.Get[da.Warehouse](factory)
	if err != nil {
		return nil, err
	}
	api, err := api.NewSimpleRest[model.Record]("/records", api.WithList(warehouse.SearchRecord))
	if err != nil {
		return nil, err
	}
	f := &Record{
		warehouse:  warehouse,
		SimpleRest: api,
	}
	return f, nil
}

// List
// @Summary      List Records with filter and page
// @Description  get records by filter
// @Tags         Warehouse
// @Accept       json
// @Produce      json
// @Param        page   query      int  true  "Page Number" minimum(1)
// @Param        page_size   query      int  true  "Page Size" minimum(10)    maximum(100)
// @Param        brand   query      string  true  "Filament Brand Name"
// @Param        color   query      string  true  "Filament Color Name"
// @Param        type   query      string  true  "Filament Type Name"
// @Success      200  {object}  api.ListResponse[model.Record]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/warehouse/records [get]
func (r *Record) List(c *gin.Context) {
	r.SimpleRest.List(c)
}
