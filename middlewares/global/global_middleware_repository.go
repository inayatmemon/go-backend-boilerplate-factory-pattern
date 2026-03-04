package global_middleware

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Repository interface {
	// RequestID adds a unique request ID to each request for tracing.
	RequestID() gin.HandlerFunc
	// GetMiddlewares returns all global middlewares to apply to every route.
	GetMiddlewares() []gin.HandlerFunc
}

type service struct {
	Input Input
}

type Input struct {
	Logger *zap.SugaredLogger
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for global middleware: %v", err)
	}
	input.Logger.Infow("Global middleware initialized")
	return &service{Input: input}
}

func (i *Input) validateInput() error {
	if i == nil {
		return errors.New("input is nil for global middleware")
	}
	if i.Logger == nil {
		return errors.New("logger is nil for global middleware")
	}
	return nil
}
