// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/dao"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/dao/org"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/service"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/server/grpc"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/server/http"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp, org.Provider))
}
