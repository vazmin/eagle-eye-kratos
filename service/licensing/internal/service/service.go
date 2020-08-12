package service

import (
	"context"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	pb "github.com/vazmin/eagle-eye-kratos/service/licensing/api"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/dao"
	"github.com/vazmin/eagle-eye-kratos/service/licensing/internal/dao/org"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.LicensingServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
	org *org.Dao
}

func (s *Service) GetLicensesByOrg(ctx context.Context, req *pb.GetLicensesByOrgReq) (*pb.Licenses, error) {

	licenses, err := s.dao.LicensesByOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return &pb.Licenses{List: licenses}, nil
}

func (s *Service) GetLicense(ctx context.Context, req *pb.GetLicenseReq) (*pb.License, error) {
	l, err := s.dao.License(ctx, req.OrganizationId, req.LicenseId)
	if err != nil {
		return nil, err
	}
	org, err := s.org.GetOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	l.Organization = org
	return l, err
}

func (s *Service) AddLicense(ctx context.Context, license *pb.License) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.AddLicense(ctx, license)
}

func (s *Service) UpdateLicense(ctx context.Context, license *pb.License) (*empty.Empty, error) {
	l, err := s.dao.License(ctx, license.OrganizationId, license.LicenseId)
	if err != nil {
		return nil, err
	}
	l.XXX_Merge(license)
	return &empty.Empty{}, s.dao.UpdateLicense(ctx, l)
}

func (s *Service) DeleteLicense(ctx context.Context, license *pb.License) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.DeleteLicense(ctx, license)
}

// New new a service and return.
func New(d dao.Dao, org *org.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
		org: org,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}

func assign(dest *pb.License, src *pb.License) {
	dest.XXX_Merge(src)
}