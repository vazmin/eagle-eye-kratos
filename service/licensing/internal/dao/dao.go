package dao

import (
	"context"
	"github.com/go-kratos/kratos/pkg/cache"
	pb "github.com/vazmin/eagle-eye-kratos/service/licensing/api"
	orgpb "github.com/vazmin/eagle-eye-kratos/service/organization/api"
	"time"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/go-kratos/kratos/pkg/time"

	"github.com/google/wire"
)

var Provider = wire.NewSet(New, NewDB, NewRedis, orgpb.NewClient)

//go:generate kratos tool genbts
// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	LicensesByOrg(ctx context.Context, orgId string) (licenses []*pb.License,err error)
	License(ctx context.Context, orgId string, licenseId string) (license *pb.License, err error)
	AddLicense(ctx context.Context, license *pb.License) (err error)
	UpdateLicense(ctx context.Context, license *pb.License) (err error)
	DeleteLicense(ctx context.Context, license *pb.License) (err error)
}

// dao dao.
type dao struct {
	db          *sql.DB
	redis       *redis.Redis
	cache *fanout.Fanout
	demoExpire int32
}


func (d *dao) LicensesByOrg(ctx context.Context, orgId string) (license []*pb.License,err error) {
	return d.RawLicensesByOrg(ctx, orgId)
}

func (d *dao) License(ctx context.Context, orgId string, licenseId string) (license *pb.License, err error) {
	addCache := true
	license, err = d.CacheLicense(ctx, licenseId)
	if err != nil {
		addCache = false
		err = nil
	}
	defer func() {
		if license != nil && license.LicenseId == "" {
			license = nil
		}
	}()
	if license != nil {
		cache.MetricHits.Inc("bts:license")
		return
	}
	cache.MetricMisses.Inc("bts:license")
	license, err = d.RawLicense(ctx, orgId, licenseId)
	if err != nil {
		return
	}
	miss := license
	if miss == nil {
		miss = &pb.License{LicenseId: ""}
	}
	if !addCache {
		return
	}
	d.cache.Do(ctx, func(c context.Context) {
		_ = d.AddCacheLicense(c, orgId, miss)
	})
	return
}



func (d *dao) AddLicense(ctx context.Context, license *pb.License) (err error) {
	err = d.RawInsertLicense(ctx, license)
	if err != nil {return err}
	err = d.AddCacheLicense(ctx, license.LicenseId, license)
	return err
}

func (d *dao) UpdateLicense(ctx context.Context, license *pb.License) (err error) {
	err = d.RawUpdateLicense(ctx, license)
	if err != nil {return err}
	err = d.AddCacheLicense(ctx, license.LicenseId, license)
	return err
}

func (d *dao) DeleteLicense(ctx context.Context, license *pb.License) (err error) {
	err = d.RawDeleteLicense(ctx, license)
	if err != nil {return err}
	err = d.DeleteCacheLicense(ctx, license.LicenseId)
	return err
}

// New new a dao and return.
func New(r *redis.Redis,  db *sql.DB) (d Dao, cf func(), err error) {
	return newDao(r, db)
}

func newDao(r *redis.Redis, db *sql.DB) (d *dao, cf func(), err error) {
	var cfg struct{
		DemoExpire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &dao{
		db: db,
		redis: r,
		cache: fanout.New("cache"),
		demoExpire: int32(time.Duration(cfg.DemoExpire) / time.Second),
	}
	cf = d.Close
	return
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return nil
}
