package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"time"
)

func NewMutex(client *redis.Client, key string, expiration time.Duration) *Mutex {
	return &Mutex{
		Key:        key,
		value:      1,
		Expiration: expiration,
		Redis:      client,
	}
}

type Mutex struct {
	Key        string
	value      uint8
	Expiration time.Duration
	Redis      *redis.Client
}

//加锁
func (m *Mutex) Lock(ctx context.Context) bool {
	if len(m.Key) == 0 {
		return false
	}
	return m.Redis.SetNX(ctx, m.Key, m.value, m.Expiration).Val()
}

//释放锁
func (m *Mutex) Unlock(ctx context.Context) int64 {
	keys := []string{m.Key}
	var script = `
		if redis.call("get",KEYS[1]) == ARGV[1] then
			return redis.call("del",KEYS[1])
		else
			return 0
		end`
	ret := m.Redis.Eval(ctx, script, keys, m.value).Val()
	return cast.ToInt64(ret)
}
