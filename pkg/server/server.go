package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	v1 "github.com/maomaozgw/filament/pkg/api/v1"
	"github.com/maomaozgw/filament/pkg/da"
	"github.com/maomaozgw/filament/pkg/model"
	"github.com/pkg/errors"
)

type Server struct {
	g      *gin.Engine
	s      *http.Server
	l      net.Listener
	cancel context.CancelFunc
}

func NewServer(opt Option) (*Server, error) {
	db, err := newOrm(opt.Orm)
	if err != nil {
		return nil, errors.Wrap(err, "new orm failed")
	}
	err = db.AutoMigrate(
		&model.Filament{},
		&model.Brand{},
		&model.Type{},
		&model.Color{},
		&model.Record{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "migrate failed")
	}
	daFactory, err := da.NewFactory(db)
	if err != nil {
		return nil, errors.Wrap(err, "new da factory failed")
	}
	router := gin.Default()
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())
	spaHandle := static.Serve("/", static.LocalFile(opt.StaticDir, true))
	router.Use(spaHandle)
	router.NoRoute(func(ctx *gin.Context) {
		ctx.File(opt.StaticDir + "/index.html")
	})
	v1Api, err := v1.New(daFactory)
	if err != nil {
		return nil, errors.Wrap(err, "new v1 api failed")
	}
	v1Api.RegisterRoutes(router.Group("/api"))
	l, err := net.Listen("tcp", opt.Addr)
	if err != nil {
		return nil, errors.Wrap(err, "listen failed")
	}
	return &Server{
		g: router,
		s: &http.Server{
			Handler: router,
		},
		l: l,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	go func() {
		select {
		case <-ctx.Done():
			shutDownCtx := context.Background()
			shutDownCtx, cancel := context.WithTimeout(shutDownCtx, 5*time.Second)
			_ = s.s.Shutdown(shutDownCtx)
			cancel()
		}
	}()
	log.Println("server started")
	return s.s.Serve(s.l)
}
