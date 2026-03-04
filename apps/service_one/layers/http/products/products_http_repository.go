package products_http_service

import (
	"errors"
	products_domain_service "go_boilerplate_project/apps/service_one/layers/domain/products"
	context_repository "go_boilerplate_project/services/context"
	api_helpers_service "go_boilerplate_project/services/helpers/api"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Repository interface {
	CreateProduct(c *gin.Context)
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
	Products products_domain_service.Repository
}

type Services struct {
	Context context_repository.Repository
}

type Helpers struct {
	API api_helpers_service.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for products http repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for products http repository")
	}
	if s.Domain == nil {
		return errors.New("domain is nil for products http repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for products http repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for products http repository")
	}
	if s.Helpers == nil {
		return errors.New("helpers is nil for products http repository")
	}
	return nil
}
