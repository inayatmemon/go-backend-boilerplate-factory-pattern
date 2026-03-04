package brands_http_service

import (
	"errors"
	brands_domain_service "go_boilerplate_project/apps/service_one/layers/domain/brands"
	context_repository "go_boilerplate_project/services/context"
	api_helpers_service "go_boilerplate_project/services/helpers/api"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Repository interface {
	EditBrand(c *gin.Context)
}

type service struct {
	Input
}

type Input struct {
	Domain   *Domain
	Logger   *zap.SugaredLogger
	Services *Services
	Helpers  *Helpers
}

type Domain struct {
	Brands brands_domain_service.Repository
}

type Services struct {
	Context context_repository.Repository
}

type Helpers struct {
	API api_helpers_service.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for brands http repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for brands http repository")
	}
	if s.Domain == nil {
		return errors.New("domain is nil for brands http repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for brands http repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for brands http repository")
	}
	if s.Helpers == nil {
		return errors.New("helpers is nil for brands http repository")
	}
	return nil
}
