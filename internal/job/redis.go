package job

import (
	"context"
	"github.com/nsevenpack/logger/v2/logger"
	"github.com/redis/go-redis/v9"
	"sync"
)

var ClientRedis *redis.Client
var CtxRedis = context.Background()
var onceRedis sync.Once

func Redis(addr string) {
	onceRedis.Do(func() {
		ClientRedis = redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   0,
		})

		// Test la connexion pour éviter les panics plus tard
		if err := ClientRedis.Ping(CtxRedis).Err(); err != nil {
			logger.Ff("❌ Erreur de connexion à Redis : %v", err)
		} else {
			logger.Sf("📡 Redis connecté sur %s", addr)
		}
	})
}
