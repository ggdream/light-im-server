package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"lim/config"
)

var (
	client   *redis.Client
	rootCtx  = context.TODO()
	baseTime = time.Second * 3
)

// Init 初始化Redis连接
func Init() error {
	cfg := config.GetRedis()

	opt := redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           0,
		MinIdleConns: 8,
		MaxIdleConns: 32,
	}
	client = redis.NewClient(&opt)

	ctx, cancel := withTimeout()
	defer cancel()

	return client.Ping(ctx).Err()
}

func Close() error {
	return client.Close()
}

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(rootCtx, baseTime)
}
