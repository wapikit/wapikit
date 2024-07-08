package internal

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func NewRedisClient(url string) *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: url,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Println("Connected to Redis successfully!!!")
			return nil
		},
	})

	return redisClient
}

func PublishMessageToRedisChannel(channel, message string) error {
	ctx := context.Background()
	err := redisClient.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}
	return nil
}
