package job

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var ClientRedis *redis.Client
var CtxRedis = context.Background()

func Redis(addr string) {
	ClientRedis = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}
