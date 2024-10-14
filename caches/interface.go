package caches

import "time"

type Cacher interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
}
