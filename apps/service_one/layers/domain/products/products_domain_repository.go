package products_domain_service

import (
	"errors"
	brands_data_service "go_boilerplate_project/apps/service_one/layers/data/brands"
	products_data_service "go_boilerplate_project/apps/service_one/layers/data/products"
	api_models "go_boilerplate_project/models/api"
	products_models "go_boilerplate_project/models/products"
	context_repository "go_boilerplate_project/services/context"
	"log"

	"go.uber.org/zap"
)

type Repository interface {
	CreateProduct(request *products_models.CreateProductRequest) *api_models.ApiResponse
}

type service struct {
	Input
}

type Input struct {
	Data     *Data
	Services *Services
	Logger   *zap.SugaredLogger
}

type Data struct {
	Products products_data_service.Repository
	Brands   brands_data_service.Repository
}

type Services struct {
	Context context_repository.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for products domain repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for products domain repository")
	}
	if s.Data == nil {
		return errors.New("data is nil for products domain repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for products domain repository")
	}
	if s.Services.Context == nil {
		return errors.New("context is nil for products domain repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for products domain repository")
	}
	return nil
}
