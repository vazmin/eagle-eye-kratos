package grpc

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	pb "github.com/vazmin/eagle-eye-kratos/service/organization/api"
)

// New new a grpc server.
func New(svc pb.OrganizationSvcServer) (ws *warden.Server, err error) {
	var (
		cfg warden.ServerConfig
		ct paladin.TOML
	)
	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	ws = warden.NewServer(&cfg)
	pb.RegisterOrganizationSvcServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}
