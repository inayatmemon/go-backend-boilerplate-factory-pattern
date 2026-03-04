package network_service

import (
	"context"
	"errors"
	network_models "go_boilerplate_project/models/network"
	"log"

	"go.uber.org/zap"
)

type Repository interface {
	Fetch(ctx context.Context, input *network_models.FetchInput) (*network_models.FetchOutput, error)
}

type service struct {
	Input Input
}

type Input struct {
	Logger *zap.SugaredLogger
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for network service: %v", err)
	}
	input.Logger.Infow("Network service initialized successfully")
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for network service")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for network service")
	}
	return nil
}
