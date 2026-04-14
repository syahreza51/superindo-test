package product

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetList(ctx context.Context, cacheKey string) (*ListProductsResponse, bool)
	SetList(ctx context.Context, cacheKey string, v *ListProductsResponse) error
	InvalidateLists(ctx context.Context) error
}

type redisCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewCache(rdb *redis.Client, ttl time.Duration) Cache {
	return &redisCache{rdb: rdb, ttl: ttl}
}

func (c *redisCache) GetList(ctx context.Context, cacheKey string) (*ListProductsResponse, bool) {
	b, err := c.rdb.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return nil, false
	}
	var v ListProductsResponse
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, false
	}
	return &v, true
}

func (c *redisCache) SetList(ctx context.Context, cacheKey string, v *ListProductsResponse) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, cacheKey, b, c.ttl).Err()
}

func (c *redisCache) InvalidateLists(ctx context.Context) error {
	const pattern = "products:list:*"
	var cursor uint64
	for {
		keys, next, err := c.rdb.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func ListCacheKey(rawQuery string) string {
	sum := sha256.Sum256([]byte(rawQuery))
	return "products:list:" + hex.EncodeToString(sum[:])
}

