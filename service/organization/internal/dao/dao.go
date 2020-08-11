package dao

import (
	"context"
	"github.com/go-kratos/kratos/pkg/cache"
	pb "github.com/vazmin/eagle-eye-kratos/service/organization/api"
	"time"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
	"github.com/go-kratos/kratos/pkg/sync/pipeline/fanout"
	xtime "github.com/go-kratos/kratos/pkg/time"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, NewDB, NewRedis)

//go:generate kratos tool genbts
// Dao dao interface
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	Organization(c context.Context, orgId string) (*pb.Organization, error)
	InsertOrganization(c context.Context, organization *pb.Organization) (err error)
	UpdateOrganization(c context.Context, organization *pb.Organization) (err error)
	DeleteOrganization(c context.Context, organization *pb.Organization) (err error)
}

// dao dao.
type dao struct {
	db          *sql.DB
	redis       *redis.Redis
	cache *fanout.Fanout
	demoExpire int32
}

// New new a dao and return.
func New(r *redis.Redis, db *sql.DB) (d Dao, cf func(), err error) {
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

func (d *dao) Organization(c context.Context, orgId string) (res *pb.Organization, err error) {
	addCache := true
	res, err = d.CacheOrg(c, orgId)
	if err != nil {
		addCache = false
		err = nil
	}
	defer func() {
		if res != nil && res.Id == "" {
			res = nil
		}
	}()
	if res != nil {
		cache.MetricHits.Inc("bts:Org")
		return
	}
	cache.MetricMisses.Inc("bts:Org")
	res, err = d.RawOrg(c, orgId)
	if err != nil {
		return
	}
	miss := res
	if miss == nil {
		miss = &pb.Organization{Id: ""}
	}
	if !addCache {
		return
	}
	d.cache.Do(c, func(c context.Context) {
		d.AddCacheOrg(c, orgId, miss)
	})
	return
}

func (d *dao) InsertOrganization(c context.Context, organization *pb.Organization) (err error) {
	err = d.RawInsertOrg(c, organization)
	if err != nil {return err}
	err = d.AddCacheOrg(c, organization.Id, organization)
	return err
}

func (d *dao) UpdateOrganization(c context.Context, organization *pb.Organization) (err error) {
	err = d.RawUpdateOrg(c, organization)
	if err != nil {return err}
	err = d.AddCacheOrg(c, organization.Id, organization)
	return err
}

func (d *dao) DeleteOrganization(c context.Context, organization *pb.Organization) (err error) {
	err = d.RawDeleteOrg(c, organization)
	if err != nil {return err}
	err = d.DeleteOrgCache(c, organization.Id)
	return err
}


