package redis

import (
	"context"
	"fmt"
	"notificationservice/src/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT),
		Password: config.REDIS_PASSWORD,
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: rdb}, nil
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Simple notification functions
func (r *RedisClient) StoreNotification(ctx context.Context, key string, data []byte) error {
	return r.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (r *RedisClient) GetNotification(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *RedisClient) GetAllNotificationKeys(ctx context.Context) ([]string, error) {
	var allKeys []string
	cursor := uint64(0)

	for {
		keys, newCursor, err := r.client.Scan(ctx, cursor, "notification:*", 100).Result()
		if err != nil {
			return nil, err
		}

		allKeys = append(allKeys, keys...)
		cursor = newCursor

		if cursor == 0 {
			break
		}
	}

	return allKeys, nil
}

func (r *RedisClient) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}
