package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(url string) *RedisClient {
	redisClient := redis.NewClient(&redis.Options{
		Addr: url,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Println("Connected to Redis successfully!!!")
			return nil
		},
	})
	return &RedisClient{redisClient}
}

func (client *RedisClient) CacheData(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	err := client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) GetCachedData(key string) (string, error) {
	ctx := context.Background()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (client *RedisClient) ComputeCacheKey(context, id, object string) string {
	return strings.Join([]string{context, object, id}, ":")
}

func (client *RedisClient) PublishMessageToRedisChannel(channel, message string) error {
	ctx := context.Background()
	err := client.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}
	return nil
}
