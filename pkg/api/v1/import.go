package v1

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type ImportData struct {
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

type Import struct {
	warehouse *da.Warehouse
}

func NewImport(factory *da.Factory) (api.RestAPI, error) {
	warehouse, err := da.Get[da.Warehouse](factory)
	if err != nil {
		return nil, err
	}
	i := &Import{
		warehouse: warehouse,
	}
	return api.NewRestAPI[ImportData]("", i), nil
}

// Create
// @Summary      Create Import Request
// @Description  Create Import Request to import filament/color/brand/type
// @Tags         Import
// @Accept       json
// @Produce      json
// @Param        import-data   body      ImportData  true  "Import Data"
// @Success      200  {object}  api.ItemResponse[ImportData]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/imports [post]
func (i *Import) Create(c *gin.Context) (*api.ItemResponse[ImportData], error) {
	data := &ImportData{}
	if err := c.ShouldBindJSON(data); err != nil {
		return nil, err
	}
	switch data.Kind {
	case "filament":
		var payload []model.Filament
		if err := json.Unmarshal(data.Data, &payload); err != nil {
			return nil, err
		}
		err := i.warehouse.Import(c, payload)
		if err != nil {
			return nil, err
		}
		return &api.ItemResponse[ImportData]{
			Code: 0,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported import kind: %s", data.Kind)
	}
}
