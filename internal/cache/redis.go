package cache

import (
	"context"
	"errors"
	"time"

	"github.com/JubeleeDev/url-shortener/internal/shortener"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(redis *redis.Client) *Cache {
	return &Cache{client: redis}
}

const applicationKey = "urlshortener:link:"

func (c *Cache) Get(ctx context.Context, code string) (string, error) {
	key := applicationKey + code
	url, err := c.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", shortener.ErrCacheMiss
	}

	return url, err
}

func (c *Cache) Set(ctx context.Context, code string, url string, ttl int) error {
	key := applicationKey + code
	_, err := c.client.Set(ctx, key, url, time.Duration(ttl)*time.Second).Result()
	return err
}

func (c *Cache) Delete(ctx context.Context, code string) error {
	key := applicationKey + code
	_, err := c.client.Del(ctx, key).Result()
	return err
}
