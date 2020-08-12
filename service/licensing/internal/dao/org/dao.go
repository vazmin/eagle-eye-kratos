package org

import (
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	"github.com/google/wire"
	orgpb "github.com/vazmin/eagle-eye-kratos/service/organization/api"
)

var Provider = wire.NewSet(New)

type Dao struct {
	orgClient orgpb.OrganizationSvcClient
}

func New() (d *Dao, err error) {
	var (
		cfg *warden.ClientConfig
		ct paladin.Map
	)
	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("OrgClient").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &Dao{}
	d.orgClient, err = orgpb.NewClient(cfg)
	return d, err
}
