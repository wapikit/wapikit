package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func newRedisClient(url string) *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr: url,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Println("Connected to Redis successfully!!!")
			return nil
		},
	})
	return redisClient
}

func GetRedisClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	redisClient = newRedisClient("localhost:6379")
	return redisClient

}

func CacheData(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	redisClient := GetRedisClient()
	err := redisClient.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetCachedData(key string) (string, error) {
	ctx := context.Background()
	redisClient := GetRedisClient()
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func ComputeCacheKey(context, id, object string) string {
	return strings.Join([]string{context, object, id}, ":")
}

func PublishMessageToRedisChannel(channel, message string) error {
	ctx := context.Background()
	redisClient := GetRedisClient()
	err := redisClient.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}
	return nil
}
