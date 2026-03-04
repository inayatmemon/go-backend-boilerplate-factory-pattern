package redis_helpers_service

import (
	"errors"
	redis_models "go_boilerplate_project/models/databases/redis"
	context_repository "go_boilerplate_project/services/context"
	"log"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Repository interface {
	Set(input *redis_models.SetInput) error
	SetNX(input *redis_models.SetNXInput) (*redis_models.SetNXOutput, error)
	Get(input *redis_models.GetInput) (*redis_models.GetOutput, error)
	Update(input *redis_models.UpdateInput) (*redis_models.UpdateOutput, error)
	Delete(input *redis_models.DeleteInput) (*redis_models.DeleteOutput, error)
	Exists(input *redis_models.ExistsInput) (*redis_models.ExistsOutput, error)
	Expire(input *redis_models.ExpireInput) (*redis_models.ExpireOutput, error)
}

type service struct {
	Input Input
}

type Input struct {
	Client   *Client
	Logger   *zap.SugaredLogger
	Services *Services
}

type Services struct {
	Context context_repository.Repository
}

type Client struct {
	RedisClient *redis.Client
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for redis helpers repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for redis helpers repository")
	}
	if s.Client == nil {
		return errors.New("client is nil for redis helpers repository")
	}
	if s.Client.RedisClient == nil {
		return errors.New("redis client is nil for redis helpers repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for redis helpers repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for redis helpers repository")
	}
	if s.Services.Context == nil {
		return errors.New("context service is nil for redis helpers repository")
	}
	return nil
}
