package application_middleware

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Repository interface {
	// AppVersion adds the application version to response headers for all service routes.
	AppVersion() gin.HandlerFunc
	// GetMiddlewares returns all application-specific middlewares.
	GetMiddlewares() []gin.HandlerFunc
}

type service struct {
	Input Input
}

type Input struct {
	Logger    *zap.SugaredLogger
	AppName   string
	AppVersion string
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for application middleware: %v", err)
	}
	input.Logger.Infow("Application middleware initialized",
		"appName", input.AppName,
		"appVersion", input.AppVersion,
	)
	return &service{Input: input}
}

func (i *Input) validateInput() error {
	if i == nil {
		return errors.New("input is nil for application middleware")
	}
	if i.Logger == nil {
		return errors.New("logger is nil for application middleware")
	}
	return nil
}
