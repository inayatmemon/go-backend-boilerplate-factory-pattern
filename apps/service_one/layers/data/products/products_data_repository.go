package products_data_service

import (
	"errors"
	common_models "go_boilerplate_project/models/commons"
	mysql_models "go_boilerplate_project/models/databases/mysql"
	products_models "go_boilerplate_project/models/products"
	context_repository "go_boilerplate_project/services/context"
	custom_helpers_service "go_boilerplate_project/services/helpers/custom"
	mongodb_helpers_service "go_boilerplate_project/services/helpers/db/mongodb"
	mysql_helpers_service "go_boilerplate_project/services/helpers/db/mysql"
	redis_helpers_service "go_boilerplate_project/services/helpers/db/redis"
	"log"

	"go.uber.org/zap"
)

type Repository interface {
	CreateProductMySQL(product *products_models.Product) error
	CreateProductMongoDB(product *products_models.Product) error
	GetCollectionName() string
	GetTableName() string

	EditProductBrandNameMySQL(brandName string, productId string, ctx ...common_models.IsContextPresentInput) error
	EditProductBrandNameMongoDB(brandName string, productId string, ctx ...common_models.IsContextPresentInput) error
}

type service struct {
	Input
}

type Input struct {
	Helpers  *Helpers
	Services *Services
	Logger   *zap.SugaredLogger
}

type Helpers struct {
	MySQL   mysql_helpers_service.Repository
	MongoDB mongodb_helpers_service.Repository
	Redis   redis_helpers_service.Repository
	Custom  custom_helpers_service.Repository
}

type Services struct {
	Context context_repository.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for products data repository: %v", err)
	}

	input.Logger.Infow("Initializing products data layer and auto-migrating table")
	// auto migrate products table
	err := input.Helpers.MySQL.AutoMigrate(&mysql_models.AutoMigrateInput{
		Model: &products_models.Product{},
	})
	if err != nil {
		log.Fatalf("Failed to auto migrate products table: %v", err)
	}

	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for products data repository")
	}
	if s.Helpers == nil {
		return errors.New("helpers is nil for products data repository")
	}
	if s.Helpers.MySQL == nil {
		return errors.New("mysql helpers is nil for products data repository")
	}
	if s.Helpers.MongoDB == nil {
		return errors.New("mongodb helpers is nil for products data repository")
	}
	if s.Helpers.Redis == nil {
		return errors.New("redis helpers is nil for products data repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for products data repository")
	}
	if s.Services.Context == nil {
		return errors.New("context service is nil for products data repository")
	}
	if s.Helpers.Custom == nil {
		return errors.New("custom helpers is nil for products data repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for products data repository")
	}
	return nil
}
