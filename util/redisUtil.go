package util

import (
	"fmt"
	"github.com/go-redis/redis"
	"gosignaler-cluster/signalerconst"
)

var (
	Redis *redis.Client
)

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", signalerconst.REDIS_HOST, signalerconst.REDIS_PORT),
		Password: signalerconst.REDIS_PASSWORD,
		DB:       signalerconst.REDIS_DBNAME, // use default DB
		PoolSize: signalerconst.REDIS_POOLSIZE,
	})
}
