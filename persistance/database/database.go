package database

import (
	"context"
	"time"

	"github.com/omegaatt36/gotasker/logging"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

// Nil defines redis returned nil value error.
const Nil = redis.Nil

var singleton redisClient

// Initialize init package.
func Initialize(ctx context.Context, addr, password string) {
	singleton.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	err := Redis().Ping(ctx).Err()
	if err != nil {
		logging.Panicf("connect to redis(%v) failed: %v", addr, err)
	}
}

// Redis returns redis client. It's safe to concurrent use.
func Redis() *redis.Client {
	if singleton.client == nil {
		panic("redis client is not created")
	}

	return singleton.client
}

// Lock implements lock using redis.
type Lock struct {
	client *redis.Client
	key    string
}

// NewLock creates lock.
func NewLock(c *redis.Client, key string) *Lock {
	return &Lock{client: c, key: key + ":lock"}
}

// Lock blocks and wait until locks.
func (l *Lock) Lock(ctx context.Context) {
	for {
		if l.TryLock(ctx) {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// TryLock try to lock and returns result.
func (l *Lock) TryLock(ctx context.Context) bool {
	return l.client.SetNX(ctx, l.key, "1", 30*time.Second).Val()
}

// Unlock unlocks.
func (l *Lock) Unlock(ctx context.Context) {
	l.client.Del(ctx, l.key)
}

// Locked returns whether locked.
func (l *Lock) Locked(ctx context.Context) bool {
	return l.client.Get(ctx, l.key).Val() == "1"
}
