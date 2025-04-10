package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/maomaozgw/filament/pkg/api"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/pkg/errors"
)

type v1 struct {
	factory    *da.Factory
	subRouters map[string][]api.RestAPI
}

// RegisterRoutes implements api.RestAPI.
func (v *v1) RegisterRoutes(g *gin.RouterGroup) {
	g = g.Group("/v1")
	for prefix, routers := range v.subRouters {
		group := g.Group(prefix)
		for _, router := range routers {
			router.RegisterRoutes(group)
		}
	}
}

func New(factory *da.Factory) (api.RestAPI, error) {
	routerCreators := map[string][]api.RouterCreator{
		"/warehouse": {
			NewFilament,
			NewRecords,
			NewStatistic,
		},
		"/meta-data": {
			NewColor,
			NewBrand,
			NewType,
			NewMetaData,
		},
		"/imports": {
			NewImport,
		},
	}
	subRouters := map[string][]api.RestAPI{}
	for prefix, err := range routerCreators {
		subRouters[prefix] = make([]api.RestAPI, len(err))
		for idx, creator := range err {
			router, err := creator(factory)
			if err != nil {
				return nil, errors.Wrapf(err, "create router %s[%d] failed", prefix, idx)
			}
			subRouters[prefix][idx] = router
		}
	}
	return &v1{
		factory:    factory,
		subRouters: subRouters,
	}, nil
}
