// internal/cache/redis_cache.go
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/epg-sync/epgsync/pkg/errors"
	"github.com/epg-sync/epgsync/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string, dest any) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.CacheMiss(key)
		}
		logger.Error("Failed to get cache",
			logger.Err(err),
			logger.String("key", key),
		)
		return errors.Wrap(err, errors.ErrCodeCacheMiss, "failed to get cache")
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeCacheInvalid, "failed to unmarshal value")
	}

	return nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeCacheInvalid, "failed to marshal value")
	}

	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		logger.Error("Failed to set cache",
			logger.Err(err),
			logger.String("key", key),
		)
		return errors.CacheWriteFailed(key, err)
	}

	return nil
}

func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		logger.Error("Failed to delete cache",
			logger.Err(err),
			logger.Strings("keys", keys),
		)
		return errors.Wrap(err, errors.ErrCodeCacheWriteFailed, "failed to delete cache")
	}

	return nil
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		logger.Error("Failed to check cache existence",
			logger.Err(err),
			logger.String("key", key),
		)
		return false, errors.Wrap(err, errors.ErrCodeCacheMiss, "failed to check existence")
	}

	return n > 0, nil
}

func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	err := r.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		logger.Error("Failed to set cache expiration",
			logger.Err(err),
			logger.String("key", key),
		)
		return errors.Wrap(err, errors.ErrCodeCacheWriteFailed, "failed to set expiration")
	}

	return nil
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
