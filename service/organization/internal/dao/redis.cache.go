package dao

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/gogo/protobuf/proto"
	pb "github.com/vazmin/eagle-eye-kratos/service/organization/api"

	"github.com/go-kratos/kratos/pkg/log"
)

var _ _redis

func (d *dao) CacheOrg(c context.Context, orgId string) (res *pb.Organization, err error) {
	key := keyOrg(orgId)
	b, err := redis.Bytes( d.redis.Do(c, "get", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
			return
		}
		log.Errorv(c, log.KV("CacheOrg", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	res = &pb.Organization{}
	err = proto.Unmarshal(b, res)
	if err != nil {
		log.Errorv(c, log.KV("CacheOrg", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// AddCacheOrg Set data to redis
func (d *dao) AddCacheOrg(c context.Context, orgId string, val *pb.Organization) (err error) {
	if val == nil {
		return
	}
	key := keyOrg(orgId)
	bytes, err := proto.Marshal(val)
	if err != nil {return err}
	if _, err = d.redis.Do(c, "set", key, bytes); err != nil {
		log.Errorv(c, log.KV("AddCacheOrg", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

// DeleteOrgCache delete data from redis
func (d *dao) DeleteOrgCache(c context.Context, orgId string) (err error) {
	key := keyOrg(orgId)
	if _, err = d.redis.Do(c, "del", key); err != nil {
		log.Errorv(c, log.KV("DeleteOrgCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
