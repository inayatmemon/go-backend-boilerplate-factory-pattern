package mysql

import (
	"context"
	"errors"
	"fmt"
	env_models "go_boilerplate_project/models/env"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlClientInstance    *gorm.DB
	mysqlClientInstanceErr error
	mysqlOnce              sync.Once
	mysqlReconnectAttempts = 5
	mysqlReconnectInterval = 2 * time.Second
	mysqlCancelMonitor     context.CancelFunc
)

// InitMySQL initializes the MySQL client and starts connection monitoring.
func InitMySQL(cfg *env_models.MySQL, logger *zap.SugaredLogger) (*gorm.DB, error) {
	mysqlOnce.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
		logger.Infof("Connecting to MySQL at %s:%d", cfg.Host, cfg.Port)

		client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			mysqlClientInstanceErr = err
			logger.Errorf("Failed to connect to MySQL: %v", err)
			log.Println("MySQL initial connection failed, attempting retry logic...")
			client, err = tryReconnectMySQL(dsn, mysqlReconnectAttempts, mysqlReconnectInterval, logger)
			if err != nil {
				mysqlClientInstanceErr = errors.New("failed to connect to MySQL after retries")
				mysqlClientInstance = nil
				return
			}
			mysqlClientInstance = client
		} else {
			mysqlClientInstance = client
			logger.Infof("Successfully connected to MySQL")
		}

		// Start monitoring and auto-reconnect handler
		go monitorMySQLConnection(cfg, dsn, logger)
	})

	return mysqlClientInstance, mysqlClientInstanceErr
}

// tryReconnectMySQL attempts to reconnect to MySQL with retry logic.
func tryReconnectMySQL(dsn string, maxAttempts int, interval time.Duration, logger *zap.SugaredLogger) (*gorm.DB, error) {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, pingErr := client.DB()
			if pingErr == nil && sqlDB.Ping() == nil {
				logger.Infof("Successfully connected to MySQL")
				return client, nil
			}
			lastErr = pingErr
			logger.Errorf("Failed to ping MySQL on reconnect (attempt %d/%d): %v", i+1, maxAttempts, pingErr)
		} else {
			lastErr = err
			logger.Errorf("Failed to connect to MySQL (attempt %d/%d): %v", i+1, maxAttempts, err)
		}
		time.Sleep(interval)
	}
	return nil, lastErr
}

// monitorMySQLConnection health-checks the MySQL connection, and reconnects if needed.
func monitorMySQLConnection(cfg *env_models.MySQL, dsn string, logger *zap.SugaredLogger) {
	ctx, cancel := context.WithCancel(context.Background())
	mysqlCancelMonitor = cancel
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if mysqlClientInstance == nil {
				logger.Infof("MySQL connection lost, attempting to reconnect...")
				client, err := tryReconnectMySQL(dsn, mysqlReconnectAttempts, mysqlReconnectInterval, logger)
				if err == nil {
					mysqlClientInstance = client
				}
				continue
			}
			sqlDB, err := mysqlClientInstance.DB()
			if err != nil {
				logger.Errorf("Failed to get underlying sql.DB: %v", err)
				client, recErr := tryReconnectMySQL(dsn, mysqlReconnectAttempts, mysqlReconnectInterval, logger)
				if recErr == nil {
					mysqlClientInstance = client
				}
				continue
			}
			pingErr := sqlDB.Ping()
			if pingErr != nil {
				logger.Errorf("MySQL connection lost: %v. Attempting to reconnect...", pingErr)
				client, recErr := tryReconnectMySQL(dsn, mysqlReconnectAttempts, mysqlReconnectInterval, logger)
				if recErr == nil {
					mysqlClientInstance = client
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// GetMySQLClient returns the shared MySQL client (gorm.DB).
func GetMySQLClient() *gorm.DB {
	return mysqlClientInstance
}

// CloseMySQLClient cleanly closes the MySQL connection and monitoring.
func CloseMySQLClient() error {
	if mysqlCancelMonitor != nil {
		mysqlCancelMonitor()
	}
	if mysqlClientInstance != nil {
		sqlDB, err := mysqlClientInstance.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func GetMySQLClientByDatabase(cfg *env_models.MySQL, logger *zap.SugaredLogger, database string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		database,
	)
	logger.Infof("Connecting to MySQL at %s:%d", cfg.Host, cfg.Port)

	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to MySQL: %v", err)
		log.Println("MySQL initial connection failed, attempting retry logic...")
		client, err = tryReconnectMySQL(dsn, mysqlReconnectAttempts, mysqlReconnectInterval, logger)
		if err != nil {
			logger.Errorf("Failed to reconnect to MySQL: %v", err)
			return nil, err
		}
		return client, nil
	}

	return client, nil
}
