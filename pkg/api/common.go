package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/da"
)

type RouterCreator func(*da.Factory) (RestAPI, error)

type ItemResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

type Pager struct {
	Page  int `json:"page" form:"page"`
	Size  int `json:"size" form:"size"`
	Total int `json:"total"`
}

type ListResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
	Pager   Pager  `json:"pager"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RestAPI interface {
	RegisterRoutes(g *gin.RouterGroup)
}

type restAPI[T any] struct {
	prefix    string
	implement any
}

func RestWrapper[T any](f func(c *gin.Context) (*T, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := f(c)
		if err != nil {
			c.JSON(500, &ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
			return
		}
		c.JSON(200, resp)
	}
}

// RegisterRoutes implements RestAPI.
func (r *restAPI[T]) RegisterRoutes(g *gin.RouterGroup) {
	group := g.Group(r.prefix)
	log.Printf("registering router %s", r.prefix)
	if list, ok := r.implement.(RestListAPI[T]); ok {
		log.Printf("registered list router %s", r.prefix)
		group.GET("", RestWrapper(list.List))
	}
	if get, ok := r.implement.(RestGetAPI[T]); ok {
		log.Printf("registered get router %s", r.prefix)
		group.GET(":id", RestWrapper(get.Get))
	}
	if create, ok := r.implement.(RestCreateAPI[T]); ok {
		log.Printf("registered create router %s", r.prefix)
		group.POST("", RestWrapper(create.Create))
	}
	if update, ok := r.implement.(RestUpdateAPI[T]); ok {
		log.Printf("registered update router %s", r.prefix)
		group.PUT(":id", RestWrapper(update.Update))
	}
	if delete, ok := r.implement.(RestDeleteAPI[T]); ok {
		log.Printf("registered delete router %s", r.prefix)
		group.DELETE(":id", RestWrapper(delete.Delete))
	}
}

func NewRestAPI[T any](prefix string, implement any) RestAPI {
	return &restAPI[T]{prefix: prefix, implement: implement}
}

type RestListAPI[T any] interface {
	List(c *gin.Context) (*ListResponse[T], error)
}
type RestGetAPI[T any] interface {
	Get(c *gin.Context) (*ItemResponse[T], error)
}
type RestCreateAPI[T any] interface {
	Create(c *gin.Context) (*ItemResponse[T], error)
}
type RestUpdateAPI[T any] interface {
	Update(c *gin.Context) (*ItemResponse[T], error)
}
type RestDeleteAPI[T any] interface {
	Delete(c *gin.Context) (*ItemResponse[T], error)
}

type SearchPageRequest struct {
	Pager
	Filter map[string]string `json:"filter"`
}

func NewListResponse[T any](items []T, total int64, page, pageSize int) *ListResponse[T] {
	return &ListResponse[T]{
		Code: 0,
		Data: items,
		Pager: Pager{
			Page:  page,
			Size:  pageSize,
			Total: int(total),
		},
	}
}

func NewItemResponse[T any](item *T) *ItemResponse[T] {
	return &ItemResponse[T]{
		Code: 0,
		Data: item,
	}
}

func NewSearchPageRequest(c *gin.Context) (*SearchPageRequest, error) {
	var result = &SearchPageRequest{
		Filter: map[string]string{},
	}
	if err := c.ShouldBindQuery(result); err != nil {
		return nil, err
	}
	if result.Page <= 0 {
		result.Page = 1
	}
	if result.Size <= 0 {
		result.Size = 100
	}
	for key := range c.Request.URL.Query() {
		if key == "page" || key == "size" {
			continue
		}
		val := c.Query(key)
		if len(val) == 0 {
			continue
		}
		result.Filter[key] = val
	}
	return result, nil
}

type IDRequest struct {
	ID uint `uri:"id"`
}

func NewIDRequest(c *gin.Context) (*IDRequest, error) {
	var result = &IDRequest{}
	if err := c.ShouldBindUri(result); err != nil {
		return nil, err
	}
	return result, nil
}

func NewListAPI[T any](
	searchFunc func(ctx context.Context, filter map[string]string, page, size int) ([]T, int64, error),
) func(c *gin.Context) (*ListResponse[T], error) {
	return func(c *gin.Context) (*ListResponse[T], error) {
		searchReq, err := NewSearchPageRequest(c)
		if err != nil {
			return nil, err
		}
		items, total, err := searchFunc(c, searchReq.Filter, searchReq.Page, searchReq.Size)
		if err != nil {
			return nil, err
		}
		return NewListResponse(items, total, searchReq.Page, searchReq.Size), nil
	}
}
