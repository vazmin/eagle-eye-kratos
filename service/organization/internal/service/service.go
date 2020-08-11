package service

import (
	"context"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	pb "github.com/vazmin/eagle-eye-kratos/service/organization/api"
	"github.com/vazmin/eagle-eye-kratos/service/organization/internal/dao"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.OrganizationSvcServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

func (s *Service) GetOrganization(ctx context.Context, req *pb.GetOrgReq) (resp *pb.Organization, err error) {
	return s.dao.Organization(ctx, req.OrganizationId)
}

func (s *Service) AddOrganization(ctx context.Context, req *pb.Organization) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.InsertOrganization(ctx, req)
}

func (s *Service) UpdateOrganization(ctx context.Context, req *pb.Organization) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.UpdateOrganization(ctx, req)
}

func (s *Service) DeleteOrganization(ctx context.Context, req *pb.Organization) (*empty.Empty, error){
	return &empty.Empty{}, s.dao.DeleteOrganization(ctx, req)
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}
