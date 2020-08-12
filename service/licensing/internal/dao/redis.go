package dao

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	pb "github.com/vazmin/eagle-eye-kratos/service/licensing/api"
)

func NewRedis() (r *redis.Redis, cf func(), err error) {
	var (
		cfg redis.Config
		ct paladin.Map
	)
	if err = paladin.Get("redis.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	r = redis.NewRedis(&cfg)
	cf = func(){r.Close()}
	return
}

func (d *dao) PingRedis(ctx context.Context) (err error) {
	if _, err = d.redis.Do(ctx, "SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

type _redis interface {
	CacheLicense(c context.Context, licenseId string) (*pb.License, error)
	AddCacheLicense(c context.Context, licenseId string, license *pb.License) (err error)
	DeleteCacheLicense(c context.Context, licenseId string) (err error)
}

func keyLicense(licenseId string) string {
	return fmt.Sprintf("license_%s", licenseId)
}
