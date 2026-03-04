package brands_data_service

import (
	"errors"
	brands_models "go_boilerplate_project/models/brands"
	common_models "go_boilerplate_project/models/commons"
	mysql_models "go_boilerplate_project/models/databases/mysql"
	context_repository "go_boilerplate_project/services/context"
	custom_helpers_service "go_boilerplate_project/services/helpers/custom"
	mongodb_helpers_service "go_boilerplate_project/services/helpers/db/mongodb"
	mysql_helpers_service "go_boilerplate_project/services/helpers/db/mysql"
	redis_helpers_service "go_boilerplate_project/services/helpers/db/redis"
	"log"

	"go.uber.org/zap"
)

type Repository interface {
	CreateBrandMySQL(brand *brands_models.Brand) error
	CreateBrandMongoDB(brand *brands_models.Brand) error
	GetCollectionName() string
	GetTableName() string
	GetBrandMySQL(brand *brands_models.Brand) error
	GetBrandMongoDB(brand *brands_models.Brand) error

	EditBrandMySQL(brandName string, brandId string, ctx ...common_models.IsContextPresentInput) error
	EditBrandMongoDB(brandName string, brandId string, ctx ...common_models.IsContextPresentInput) error
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
	Redis   redis_helpers_service.Repository
	MongoDB mongodb_helpers_service.Repository
	Custom  custom_helpers_service.Repository
}

type Services struct {
	Context context_repository.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for brands data repository: %v", err)
	}

	input.Logger.Infow("Initializing brands data layer and auto-migrating table")
	// auto migrate brands table
	err := input.Helpers.MySQL.AutoMigrate(&mysql_models.AutoMigrateInput{
		Model: &brands_models.Brand{},
	})
	if err != nil {
		log.Fatalf("Failed to auto migrate brands table: %v", err)
	}

	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for brands data repository")
	}
	if s.Helpers == nil {
		return errors.New("helpers is nil for brands data repository")
	}
	if s.Helpers.MySQL == nil {
		return errors.New("mysql helpers is nil for brands data repository")
	}
	if s.Helpers.MongoDB == nil {
		return errors.New("mongodb helpers is nil for brands data repository")
	}
	if s.Helpers.Redis == nil {
		return errors.New("redis helpers is nil for brands data repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for brands data repository")
	}
	if s.Services.Context == nil {
		return errors.New("context service is nil for brands data repository")
	}
	if s.Helpers.Custom == nil {
		return errors.New("custom helpers is nil for brands data repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for brands data repository")
	}
	return nil
}
