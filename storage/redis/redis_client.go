package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	env_models "go_boilerplate_project/models/env"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	clientInstance      *redis.Client
	clientInstanceErr   error
	redisOnce           sync.Once
	reconnectAttempts   = 5
	reconnectInterval   = 2 * time.Second
	backgroundCancelFun context.CancelFunc
)

// InitRedis initializes the Redis client, and starts connection monitoring and auto-reconnect.
func InitRedis(cfg *env_models.Redis, logger *zap.SugaredLogger) (*redis.Client, error) {
	redisOnce.Do(func() {
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		logger.Infof("Connecting to Redis at %s", addr)
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Password,
			DB:       cfg.Database,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := client.Ping(ctx).Err()
		if err != nil {
			log.Printf("Redis initial ping failed: %v, attempting retry logic...", err)
			client = tryReconnectRedis(cfg, reconnectAttempts, reconnectInterval, logger)
			if client == nil {
				clientInstanceErr = errors.New("failed to connect to Redis after retries")
				clientInstance = nil
				return
			}
			clientInstance = client
		} else {
			clientInstance = client
			logger.Infof("Successfully connected to Redis")
		}

		// Start monitoring and auto-reconnect handler
		go monitorRedisConnection(cfg, logger)
	})

	return clientInstance, clientInstanceErr
}

// tryReconnectRedis attempts to reconnect to Redis with retry logic.
func tryReconnectRedis(cfg *env_models.Redis, maxAttempts int, interval time.Duration, logger *zap.SugaredLogger) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	var lastClient *redis.Client
	for i := 0; i < maxAttempts; i++ {
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Password,
			DB:       cfg.Database,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := client.Ping(ctx).Err()
		cancel()
		if err == nil {
			logger.Infof("Successfully connected to Redis")
			return client
		}

		lastClient = client
		logger.Errorf("Failed to connect to Redis (attempt %d/%d): %v", i+1, maxAttempts, err)
		time.Sleep(interval)
	}
	return lastClient
}

// monitorRedisConnection health-checks the Redis connection, and reconnects if needed.
func monitorRedisConnection(cfg *env_models.Redis, logger *zap.SugaredLogger) {
	ctx, cancel := context.WithCancel(context.Background())
	backgroundCancelFun = cancel
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if clientInstance == nil {
				logger.Infof("Redis connection lost, attempting to reconnect...")
				clientInstance = tryReconnectRedis(cfg, reconnectAttempts, reconnectInterval, logger)
				continue
			}
			pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)
			err := clientInstance.Ping(pingCtx).Err()
			pingCancel()
			if err != nil {
				logger.Errorf("Redis connection lost: %v. Attempting to reconnect...", err)
				clientInstance = tryReconnectRedis(cfg, reconnectAttempts, reconnectInterval, logger)
			}
		case <-ctx.Done():
			return
		}
	}
}

// GetRedisClient returns the shared Redis client instance.
func GetRedisClient() *redis.Client {
	return clientInstance
}

// CloseRedisClient cleanly closes the Redis connection and monitoring.
func CloseRedisClient() error {
	if backgroundCancelFun != nil {
		backgroundCancelFun()
	}
	if clientInstance != nil {
		return clientInstance.Close()
	}
	return nil
}
