package serviceone_router

import (
	"errors"
	brands_http_service "go_boilerplate_project/apps/service_one/layers/http/brands"
	products_http_service "go_boilerplate_project/apps/service_one/layers/http/products"
	env_models "go_boilerplate_project/models/env"
	context_repository "go_boilerplate_project/services/context"
	api_helpers_service "go_boilerplate_project/services/helpers/api"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Repository interface {
	ConfigureRouter()
	SetupRoutes()
	Run()
}

type service struct {
	Input
	Engine
}

type Input struct {
	Logger   *zap.SugaredLogger
	Services *Services
	Helpers  *Helpers
	Env      *env_models.Environment
	Http     *Http
}

type Services struct {
	Context context_repository.Repository
}

type Helpers struct {
	API api_helpers_service.Repository
}

type Engine struct {
	Engine *gin.Engine
}

type Http struct {
	Brands   brands_http_service.Repository
	Products products_http_service.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for serviceone router repository: %v", err)
	}
	engine := Engine{
		Engine: gin.New(),
	}
	return &service{
		Input:  input,
		Engine: engine,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for serviceone router repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for serviceone router repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for serviceone router repository")
	}
	if s.Helpers == nil {
		return errors.New("helpers is nil for serviceone router repository")
	}

	if s.Services.Context == nil {
		return errors.New("context is nil for serviceone router repository")
	}
	if s.Helpers.API == nil {
		return errors.New("api helpers is nil for serviceone router repository")
	}
	if s.Env == nil {
		return errors.New("env is nil for serviceone router repository")
	}
	if s.Http == nil {
		return errors.New("http is nil for serviceone router repository")
	}
	if s.Http.Brands == nil {
		return errors.New("brands http is nil for serviceone router repository")
	}
	if s.Http.Products == nil {
		return errors.New("products http is nil for serviceone router repository")
	}
	return nil
}
