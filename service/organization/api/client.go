package api

import (
	"context"
	"github.com/go-kratos/kratos/pkg/naming/etcd"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden/resolver"

	"google.golang.org/grpc"
)

func init()  {
	resolver.Register(etcd.Builder(nil))
}
// AppID .
//const AppID = "127.0.0.1:9000"
//var target = fmt.Sprintf("direct://default/%s", AppID)
const AppID = "organizations.service"
const target = "etcd://default/" + AppID
// NewClient new grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (OrganizationSvcClient, error) {
	client := warden.NewClient(cfg, opts...)
	cc, err := client.Dial(context.Background(), target)
	if err != nil {
		return nil, err
	}
	return NewOrganizationSvcClient(cc), nil
}

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc --bm api.proto
