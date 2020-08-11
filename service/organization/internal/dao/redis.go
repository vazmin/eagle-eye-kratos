package dao

import (
	"context"
	"fmt"
	"github.com/vazmin/eagle-eye-kratos/service/organization/internal/model"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
)

type _redis interface {
	// mc: -key=keyArt -type=get
	CacheOrg(c context.Context, orgId string) (*model.Article, error)
	// mc: -key=keyArt -expire=d.demoExpire
	AddCacheOrg(c context.Context, orgId string, art *model.Article) (err error)
	// mc: -key=keyArt
	DeleteOrgCache(c context.Context, orgId string) (err error)
}

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

func keyOrg(orgId string) string {
	return fmt.Sprintf("org_%s", orgId)
}