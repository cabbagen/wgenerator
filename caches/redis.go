package caches

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisCacher Cacher = nil

type RedisCacher struct {
	ctx    context.Context
	cacher *redis.Client
}

func (rc RedisCacher) Get(key string) (interface{}, error) {
	return rc.cacher.Get(rc.ctx, key).Result()
}

func (rc RedisCacher) Set(key string, value interface{}, expiration time.Duration) error {
	return rc.cacher.Set(rc.ctx, key, value, expiration).Err()
}

func InitRedisCacherInstance(db int, address, password string) {
	if redisCacher != nil {
		return
	}
	redisCacher = RedisCacher{
		ctx:    context.Background(),
		cacher: redis.NewClient(&redis.Options{DB: db, Addr: address, Password: password}),
	}
}

func GetRedisCacherInstance() Cacher {
	if redisCacher == nil {
		fmt.Printf("you should connect the redis before invoke `GetRedisCacherInstance` method\n\n")
		return nil
	}
	return redisCacher
}
