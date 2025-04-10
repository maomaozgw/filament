package v1

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
)

type Filament struct {
	warehouse *da.Warehouse
}

func NewFilament(factory *da.Factory) (api.RestAPI, error) {
	warehouse, err := da.Get[da.Warehouse](factory)
	if err != nil {
		return nil, err
	}
	f := &Filament{
		warehouse: warehouse,
	}
	return api.NewRestAPI[model.Filament]("/filaments", f), nil
}

// List
// @Summary      List Filament with filter and pagen
// @Description  get filaments by filter
// @Tags         Warehouse
// @Accept       json
// @Produce      json
// @Param        page   query      int  true  "Page Number" minimum(1)
// @Param        page_size   query      int  true  "Page Size" minimum(10)    maximum(100)
// @Param        brand   query      string  true  "Filament Brand Name"
// @Param        color   query      string  true  "Filament Color Name"
// @Param        type   query      string  true  "Filament Type Name"
// @Success      200  {object}  api.ListResponse[model.Filament]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/warehouse/filaments [get]
func (f *Filament) List(c *gin.Context) (*api.ListResponse[model.Filament], error) {
	req, err := api.NewSearchPageRequest(c)
	if err != nil {
		return nil, err
	}
	items, total, err := f.warehouse.SearchWarehouse(c, req.Filter, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &api.ListResponse[model.Filament]{
		Code:    0,
		Message: "success",
		Data:    items,
		Pager: api.Pager{
			Page:  1,
			Size:  100,
			Total: int(total),
		},
	}, nil
}

// Get
// @Summary      Get Filament by ID
// @Description  get filament by ID
// @Tags         Warehouse
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Filament ID"
// @Success      200  {object}  api.ItemResponse[model.Filament]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      404  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/warehouse/filaments/{id} [get]
func (f *Filament) Get(c *gin.Context) (*api.ItemResponse[model.Filament], error) {
	id := &api.IDRequest{}
	if err := c.ShouldBindQuery(id); err != nil {
		return nil, err
	}
	item, err := f.warehouse.GetFilament(c, id.ID)
	if err != nil {
		return nil, err
	}
	return &api.ItemResponse[model.Filament]{
		Code:    0,
		Message: "success",
		Data:    item,
	}, nil
}

// Create
// @Summary      Stock in Filament
// @Description  stock in filament, add/increase filament to warehouse
// @Tags         Warehouse
// @Accept       json
// @Produce      json
// @Param        filament   body      model.Filament  true  "Filament"
// @Success      200  {object}  api.ItemResponse[model.Filament]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/warehouse/filaments [post]
func (f *Filament) Create(c *gin.Context) (*api.ItemResponse[model.Filament], error) {
	record := &model.Filament{}
	if err := c.ShouldBindJSON(record); err != nil {
		return nil, err
	}
	if record.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if record.Price < 0 {
		return nil, errors.New("price must be greater or equal 0")
	}
	if err := f.warehouse.StockIn(c, record); err != nil {
		return nil, err
	}
	return &api.ItemResponse[model.Filament]{
		Code:    0,
		Message: "success",
		Data:    &model.Filament{},
	}, nil
}

// Update godoc
// @Summary
// @Description  Stock take Filament
// @Description  stock take filament, update filament quantity in warehouse
// @Tags         Warehouse
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Filament ID"
// @Param        filament   body      model.Filament  true  "Filament"
// @Success      200  {object}  api.ItemResponse[model.Filament]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      404  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/warehouse/filaments/{id} [put]
func (f *Filament) Update(c *gin.Context) (*api.ItemResponse[model.Filament], error) {
	id := &api.IDRequest{}
	if err := c.ShouldBindUri(id); err != nil {
		return nil, err
	}
	newItem := &model.Filament{}
	if err := c.ShouldBindJSON(newItem); err != nil {
		return nil, err
	}
	if newItem.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	newItem.ID = id.ID
	err := f.warehouse.StockTake(c, newItem)
	if err != nil {
		return nil, err
	}
	return &api.ItemResponse[model.Filament]{
		Code:    0,
		Message: "success",
		Data:    &model.Filament{},
	}, nil
}

// Delete
// @Summary      Stock out Filament
// @Description  stock out filament, decrease filament from warehouse
// @Tags         Warehouse
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Filament ID"
// @Param        filament   body      model.Filament  true  "Filament"
// @Success      200  {object}  api.ItemResponse[model.Filament]
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /api/v1/warehouse/filaments/{id} [delete]
func (f *Filament) Delete(c *gin.Context) (*api.ItemResponse[model.Filament], error) {
	idReq := &api.IDRequest{}
	if err := c.ShouldBindUri(idReq); err != nil {
		return nil, err
	}
	filament := &model.Filament{}
	if err := c.ShouldBindJSON(filament); err != nil {
		return nil, err
	}
	if filament.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	filament.ID = idReq.ID
	log.Printf("%+v", filament)
	if err := f.warehouse.StockOut(c, filament); err != nil {
		return nil, err
	}
	return &api.ItemResponse[model.Filament]{
		Code:    0,
		Message: "success",
		Data:    &model.Filament{},
	}, nil
}
