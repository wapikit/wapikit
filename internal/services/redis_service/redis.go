package cache_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	redisPoolLib "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	RedSync         *redsync.Redsync
	RateLimitPrefix string
}

func NewRedisClient(url string) *RedisClient {
	fmt.Println("Connecting to Redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr: url,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Println("Connected to Redis successfully!!!")
			return nil
		},
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis: ", err)
		return nil
	}

	pool := redisPoolLib.NewPool(redisClient)
	redSync := redsync.New(pool)

	return &RedisClient{
		redisClient,
		redSync,
		"wapikit:rate_limit",
	}
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

func (client *RedisClient) PublishMessageToRedisChannel(channel string, message []byte) error {
	fmt.Println("Publishing message to Redis channel...")
	ctx := context.Background()
	err := client.Publish(ctx, channel, message).Err()
	if err != nil {
		fmt.Println("Error publishing message to Redis channel: ", err.Error())
		return err
	}
	fmt.Println("Message published to Redis channel successfully!!!")
	return nil
}

func (client *RedisClient) ComputeRateLimitKey(ipAddress, path string) string {
	return strings.Join([]string{client.RateLimitPrefix, ipAddress, path}, ":")
}

func (client *RedisClient) AcquireLock(lockKey string, ttl time.Duration) (*redsync.Mutex, error) {
	// Create a mutex with the specified key and TTL
	mutex := client.RedSync.NewMutex(lockKey, redsync.WithExpiry(ttl))

	// Try to acquire the lock
	if err := mutex.Lock(); err != nil {
		return nil, err
	}
	return mutex, nil
}

func (client *RedisClient) ReleaseLock(mutex *redsync.Mutex) error {
	if _, err := mutex.Unlock(); err != nil {
		return err
	}
	return nil
}
