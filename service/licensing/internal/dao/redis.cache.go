package dao

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gogo/protobuf/proto"
	pb "github.com/vazmin/eagle-eye-kratos/service/licensing/api"
)

var _ _redis

func (d *dao) CacheLicense(c context.Context, licenseId string) (res *pb.License, err error) {
	key := keyLicense(licenseId)
	b, err := redis.Bytes( d.redis.Do(c, "get", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
			return
		}
		log.Errorv(c, log.KV("CacheLicense", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	res = &pb.License{}
	err = proto.Unmarshal(b, res)
	if err != nil {
		log.Errorv(c, log.KV("CacheLicense", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

func (d *dao) AddCacheLicense(c context.Context, licenseId string, license *pb.License) (err error) {
	if license == nil {
		return
	}
	key := keyLicense(licenseId)
	bytes, err := proto.Marshal(license)
	if err != nil {return err}
	if _, err = d.redis.Do(c, "set", key, bytes); err != nil {
		log.Errorv(c, log.KV("AddCacheLicense", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}

func (d *dao) DeleteCacheLicense(c context.Context, licenseId string) (err error) {
	key := keyLicense(licenseId)
	if _, err = d.redis.Do(c, "del", key); err != nil {
		log.Errorv(c, log.KV("DeleteOrgCache", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
