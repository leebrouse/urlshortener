package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"githum.com/leebrouse/urlshortener/config"
	"githum.com/leebrouse/urlshortener/internal/repo"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg config.RedisConfig) (*RedisCache,error) {
	client:=redis.NewClient(&redis.Options{
		Addr: cfg.Address,
		Password: cfg.Password,
		DB: cfg.DB,
	})

	if err:=client.Ping(context.Background()).Err();err!=nil{
		return nil,err
	}

	return &RedisCache{client: client},nil
}

func (c *RedisCache) SetURL(ctx context.Context, url repo.Url) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}
	if err := c.client.Set(ctx, url.ShortUrl, data, time.Until(url.ExpiresAt)).Err(); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) GetURL(ctx context.Context, shortCode string) (*repo.Url, error) {
	data, err := c.client.Get(ctx, shortCode).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var url repo.Url
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, err
	}

	return &url, nil

}

func (c *RedisCache) Close() error {
	return c.client.Close()
}