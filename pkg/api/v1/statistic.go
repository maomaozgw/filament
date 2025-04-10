package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type Statistic struct {
	warehouse *da.Warehouse
	*api.SimpleRest[model.Statistic]
}

func NewStatistic(factory *da.Factory) (api.RestAPI, error) {
	warehouse, err := da.Get[da.Warehouse](factory)
	if err != nil {
		return nil, err
	}
	restAPi, err := api.NewSimpleRest(
		"/statistic",
		api.WithList(warehouse.SearchStatistic),
	)
	if err != nil {
		return nil, err
	}
	w := &Statistic{
		warehouse:  warehouse,
		SimpleRest: restAPi,
	}
	return w, nil
}

// List
// @Summary Get all statistics
// @Description Get all statistics
// @Tags Warehouse
// @Accept json
// @Produce json
// @Success 200 {object} api.ListResponse[model.Statistic]
// @Router /v1/warehouse/statistic [get]
func (s *Statistic) List(c *gin.Context) {
	s.SimpleRest.List(c)
}
