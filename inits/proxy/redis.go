package proxy

import (
	"github.com/lfxnxf/zdy_tools/inits"
	"github.com/lfxnxf/zdy_tools/resource/redis"
)

type Redis struct {
	*redis.Redis
}

func InitRedis(name string) *Redis {
	return &Redis{inits.RedisClient(name)}
}
