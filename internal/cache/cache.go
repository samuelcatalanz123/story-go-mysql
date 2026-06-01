// Package cache is a small wrapper over Redis used for two things: caching
// list results (cache-aside) and rate limiting. It degrades gracefully: if
// Redis is unavailable the app uses Noop and keeps working (just slower, with
// no caching or limiting).
package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache is the minimal interface the app needs.
type Cache interface {
	// Get returns the cached bytes and whether the key was found.
	Get(ctx context.Context, key string) ([]byte, bool, error)
	// Set stores value under key with a time-to-live (expiration).
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// DelByPrefix deletes every key starting with prefix (uses SCAN, not KEYS,
	// so it never blocks Redis).
	DelByPrefix(ctx context.Context, prefix string) error
	// Incr increments a counter, setting its expiry on the first hit, and
	// returns the new value. Used for rate limiting.
	Incr(ctx context.Context, key string, window time.Duration) (int64, error)
}

// --- Noop: caché desactivada (cuando no hay Redis) ---

// Noop does nothing: Get never finds, Set/Del are ignored, Incr returns 0 so
// rate limiting is effectively disabled. Lets the app run without Redis.
type Noop struct{}

func (Noop) Get(context.Context, string) ([]byte, bool, error)          { return nil, false, nil }
func (Noop) Set(context.Context, string, []byte, time.Duration) error   { return nil }
func (Noop) DelByPrefix(context.Context, string) error                  { return nil }
func (Noop) Incr(context.Context, string, time.Duration) (int64, error) { return 0, nil }

// --- Redis ---

// Redis implements Cache backed by a Redis server.
type Redis struct {
	client *redis.Client
}

// NewRedis connects to Redis at addr and verifies the connection with PING.
func NewRedis(ctx context.Context, addr, password string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: password})
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &Redis{client: client}, nil
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil // key not present
	}
	if err != nil {
		return nil, false, err
	}
	return val, true, nil
}

func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) DelByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	for {
		keys, next, err := r.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (r *Redis) Incr(ctx context.Context, key string, window time.Duration) (int64, error) {
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		// First hit in this window: set the expiry.
		if err := r.client.Expire(ctx, key, window).Err(); err != nil {
			return count, err
		}
	}
	return count, nil
}
