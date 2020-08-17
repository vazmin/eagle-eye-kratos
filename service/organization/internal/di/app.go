package di

import (
	"context"
	"github.com/go-kratos/kratos/pkg/conf/env"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/naming"
	"github.com/go-kratos/kratos/pkg/naming/etcd"
	"github.com/go-kratos/kratos/pkg/net/trace/zipkin"
	"github.com/vazmin/eagle-eye-kratos/service/organization/api"
	"time"

	"github.com/vazmin/eagle-eye-kratos/service/organization/internal/service"
	appenv "github.com/vazmin/eagle-eye-kratos/common/env"

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
	InitZipkin()
	cf := DiscoveryRegister()
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err := g.Shutdown(ctx); err != nil {
			log.Error("grpcSrv.Shutdown error(%v)", err)
		}
		if err := h.Shutdown(ctx); err != nil {
			log.Error("httpSrv.Shutdown error(%v)", err)
		}
		cf()
		cancel()
	}
	return
}

func DiscoveryRegister() (closeFunc func()) {
	var (
		cfg warden.ServerConfig
		ct paladin.TOML
	)
	if err := paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err := ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	//hn, _ := os.Hostname()
	dis, err := etcd.New(nil)
	if err != nil {
		panic(err)
	}
	ins := &naming.Instance{
		Zone:     env.Zone,
		Env:      env.DeployEnv,
		AppID:    api.AppID,
		//Hostname: hn,
		Addrs: []string{
			"grpc://" + cfg.Addr,
		},
	}
	cancel, err := dis.Register(context.Background(), ins)
	if err != nil {
		panic(err)
	}

	return cancel
}

func InitZipkin()  {
	ep := appenv.ZipkinEndpoint
	if ep != "" {
		zipkin.Init(&zipkin.Config{
			Endpoint: ep + "/api/v2/spans",
		})
	}
}