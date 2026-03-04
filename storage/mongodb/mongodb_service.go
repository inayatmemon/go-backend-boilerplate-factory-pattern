package mongodb

import (
	"context"
	"errors"
	"fmt"
	env_models "go_boilerplate_project/models/env"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
	reconnectAttempts = 5
	reconnectInterval = 2 * time.Second
	backgroundCancel  context.CancelFunc
)

// InitMongoDB initializes the MongoDB client and starts connection monitoring.
func InitMongoDB(cfg *env_models.MongoDB, logger *zap.SugaredLogger) (*mongo.Client, error) {
	mongoOnce.Do(func() {
		uri := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
		logger.Infof("Connecting to MongoDB at %s", uri)
		clientOptions := options.Client().ApplyURI(uri)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientInstanceErr = err
			logger.Errorf("Failed to connect to MongoDB: %v", err)
			return
		}

		// Try ping
		if err := client.Ping(ctx, nil); err != nil {
			clientInstanceErr = err
			clientInstance = nil
			log.Println("MongoDB initial ping failed, attempting retry logic...")
			clientInstance = tryReconnect(uri, reconnectAttempts, reconnectInterval, logger)
			if clientInstance == nil {
				clientInstanceErr = errors.New("failed to connect to MongoDB after retries")
				return
			}
		} else {
			clientInstance = client
			logger.Infof("Successfully connected to MongoDB")
		}

		// Start monitoring and auto-reconnect handler
		go monitorConnection(uri, logger)
	})

	return clientInstance, clientInstanceErr
}

// tryReconnect attempts to reconnect to MongoDB with retry logic.
func tryReconnect(uri string, maxAttempts int, interval time.Duration, logger *zap.SugaredLogger) *mongo.Client {
	var lastClient *mongo.Client
	for i := 0; i < maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		pingErr := error(nil)
		if err == nil {
			pingErr = client.Ping(ctx, nil)
			if pingErr != nil {
				logger.Errorf("Failed to ping MongoDB: %v", pingErr)
			}
		}
		cancel()
		if err == nil && pingErr == nil {
			logger.Infof("Successfully connected to MongoDB")
			return client
		}
		lastClient = client
		logger.Errorf("Failed to connect to MongoDB: %v", err)
		logger.Infof("Attempting to reconnect to MongoDB... %d/%d", i+1, maxAttempts)
		time.Sleep(interval)
	}
	return lastClient
}

// monitorConnection health-checks the MongoDB connection, and reconnects if needed.
func monitorConnection(uri string, logger *zap.SugaredLogger) {
	ctx, cancel := context.WithCancel(context.Background())
	backgroundCancel = cancel
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if clientInstance == nil {
				logger.Infof("MongoDB connection lost, attempting to reconnect...")
				clientInstance = tryReconnect(uri, reconnectAttempts, reconnectInterval, logger)
				continue
			}
			pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
			err := clientInstance.Ping(pingCtx, nil)
			pingCancel()
			if err != nil {
				logger.Errorf("MongoDB connection lost: %v. Attempting to reconnect...", err)
				clientInstance = tryReconnect(uri, reconnectAttempts, reconnectInterval, logger)
			}
		case <-ctx.Done():
			return
		}
	}
}

// GetMongoClient returns the shared MongoDB client.
func GetMongoClient() *mongo.Client {
	return clientInstance
}

// CloseMongoClient cleanly closes the MongoDB connection and monitoring.
func CloseMongoClient() error {
	if backgroundCancel != nil {
		backgroundCancel()
	}
	if clientInstance != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return clientInstance.Disconnect(ctx)
	}
	return nil
}
