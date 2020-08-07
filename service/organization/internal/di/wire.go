// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"github.com/vazmin/eagle-eye-kratos/service/organization/internal/dao"
	"github.com/vazmin/eagle-eye-kratos/service/organization/internal/service"
	"github.com/vazmin/eagle-eye-kratos/service/organization/internal/server/grpc"
	"github.com/vazmin/eagle-eye-kratos/service/organization/internal/server/http"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp))
}
