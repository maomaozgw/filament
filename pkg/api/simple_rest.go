package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type restOption[T any] struct {
	listFunc   func(ctx context.Context, filter map[string]string, page, pageSize int) ([]T, int64, error)
	getFunc    func(ctx context.Context, id uint, name string) (*T, error)
	createFunc func(ctx context.Context, item *T) (*T, error)
	updateFunc func(ctx context.Context, id uint, item *T) (*T, error)
	deleteFunc func(ctx context.Context, id uint) (*T, error)
}

type RestOpt[T any] func(*restOption[T]) error

func WithList[T any](fn func(ctx context.Context, filter map[string]string, page, pageSize int) ([]T, int64, error)) RestOpt[T] {
	return func(opt *restOption[T]) error {
		opt.listFunc = fn
		return nil
	}
}

func WithGet[T any](fn func(ctx context.Context, id uint, name string) (*T, error)) RestOpt[T] {
	return func(opt *restOption[T]) error {
		opt.getFunc = fn
		return nil
	}
}

func WithCreate[T any](fn func(ctx context.Context, item *T) (*T, error)) RestOpt[T] {
	return func(opt *restOption[T]) error {
		opt.createFunc = fn
		return nil
	}
}
func WithUpdate[T any](fn func(ctx context.Context, id uint, item *T) (*T, error)) RestOpt[T] {
	return func(opt *restOption[T]) error {
		opt.updateFunc = fn
		return nil
	}
}

func WithDelete[T any](fn func(ctx context.Context, id uint) (*T, error)) RestOpt[T] {
	return func(opt *restOption[T]) error {
		opt.deleteFunc = fn
		return nil
	}
}

type SimpleRest[T any] struct {
	opt    restOption[T]
	prefix string
}

func NewSimpleRest[T any](prefix string, opts ...RestOpt[T]) (*SimpleRest[T], error) {
	var opt = restOption[T]{}
	for _, o := range opts {
		if err := o(&opt); err != nil {
			return nil, err
		}
	}
	return &SimpleRest[T]{
		opt:    opt,
		prefix: prefix,
	}, nil
}

func (s *SimpleRest[T]) RegisterRoutes(g *gin.RouterGroup) {
	g = g.Group(s.prefix)
	if s.opt.listFunc != nil {
		g.GET("", s.List)
	}
	if s.opt.getFunc != nil {
		g.GET(":id", s.Get)
	}
	if s.opt.createFunc != nil {
		g.POST("", s.Create)
	}
	if s.opt.updateFunc != nil {
		g.PUT(":id", s.Update)
	}
	if s.opt.deleteFunc != nil {
		g.DELETE(":id", s.Delete)
	}
}

func (s *SimpleRest[T]) List(c *gin.Context) {
	req, err := NewSearchPageRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	items, total, err := s.opt.listFunc(c, req.Filter, req.Page, req.Size)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, &ListResponse[T]{
		Code:    0,
		Message: "success",
		Data:    items,
		Pager: Pager{
			Page:  req.Page,
			Size:  req.Size,
			Total: int(total),
		},
	})
}

func (s *SimpleRest[T]) Create(c *gin.Context) {
	var item = new(T)
	if err := c.ShouldBindJSON(&item); err != nil {
		_ = c.Error(err)
		return
	}
	item, err := s.opt.createFunc(c, item)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, &ItemResponse[T]{
		Code:    0,
		Message: "success",
		Data:    item,
	})
}

func (s *SimpleRest[T]) Get(c *gin.Context) {
	idReq, err := NewIDRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	item, err := s.opt.getFunc(c, idReq.ID, "")
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, &ItemResponse[T]{
		Code:    0,
		Message: "success",
		Data:    item,
	})
}

func (s *SimpleRest[T]) Update(c *gin.Context) {
	idReq, err := NewIDRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var item = new(T)
	if err := c.ShouldBindJSON(&item); err != nil {
		_ = c.Error(err)
		return
	}
	item, err = s.opt.updateFunc(c, idReq.ID, item)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, &ItemResponse[T]{
		Code:    0,
		Message: "success",
		Data:    item,
	})
}

func (s *SimpleRest[T]) Delete(c *gin.Context) {
	idReq, err := NewIDRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	item, err := s.opt.deleteFunc(c, idReq.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, &ItemResponse[T]{
		Code:    0,
		Message: "success",
		Data:    item,
	})
}
