package di

import (
	"context"
	"github.com/go-kratos/kratos/pkg/net/trace/zipkin"
	appenv "github.com/vazmin/eagle-eye-kratos/common/env"
	"time"

	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/service"

	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
)

//go:generate kratos tool wire
type App struct {
	svc *service.Service
	http *bm.Engine
	grpc *warden.Server
}

func NewApp(svc *service.Service, h *bm.Engine, g *warden.Server) (app *App, closeFunc func(), err error){
	app = &App{
		svc: svc,
		http: h,
		grpc: g,
	}
	appenv.Init()
	InitZipkin()
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err := g.Shutdown(ctx); err != nil {
			log.Error("grpcSrv.Shutdown error(%v)", err)
		}
		if err := h.Shutdown(ctx); err != nil {
			log.Error("httpSrv.Shutdown error(%v)", err)
		}
		cancel()
	}
	return
}

func InitZipkin()  {
	ep := appenv.ZipkinEndpoint
	if ep != "" {
		zipkin.Init(&zipkin.Config{
			Endpoint: ep + "/api/v2/spans",
		})
	}
}