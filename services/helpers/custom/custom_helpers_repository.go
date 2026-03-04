package custom_helpers_service

import (
	"errors"
	common_models "go_boilerplate_project/models/commons"
	"log"

	"go.uber.org/zap"
)

type Repository interface {
	IsContextPresent(input ...common_models.IsContextPresentInput) bool
}

type service struct {
	Input
}

type Input struct {
	Logger *zap.SugaredLogger
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for custom helpers repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for custom helpers repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for custom helpers repository")
	}
	return nil
}
