package mysql_helpers_service

import (
	"errors"
	mysql_models "go_boilerplate_project/models/databases/mysql"
	env_models "go_boilerplate_project/models/env"
	context_repository "go_boilerplate_project/services/context"
	"log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	Create(input *mysql_models.CreateInput) error
	CreateInBatches(input *mysql_models.CreateInBatchesInput) error

	FindOne(input *mysql_models.FindOneInput) error
	FindMany(input *mysql_models.FindManyInput) error

	Update(input *mysql_models.UpdateInput) (*mysql_models.UpdateOutput, error)
	Delete(input *mysql_models.DeleteInput) (*mysql_models.DeleteOutput, error)

	RawQuery(input *mysql_models.RawQueryInput) error
	Exec(input *mysql_models.ExecInput) (*mysql_models.ExecOutput, error)

	RunTransaction(input *mysql_models.TransactionInput) error

	AutoMigrate(input *mysql_models.AutoMigrateInput) error
}

type service struct {
	Input Input
}

type Input struct {
	Client   *Client
	Logger   *zap.SugaredLogger
	Services *Services
	Env      *env_models.MySQL
}

type Services struct {
	Context context_repository.Repository
}

type Client struct {
	MySQLClient *gorm.DB
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for mysql helpers repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for mysql helpers repository")
	}
	if s.Client == nil {
		return errors.New("client is nil for mysql helpers repository")
	}
	if s.Client.MySQLClient == nil {
		return errors.New("mysql client is nil for mysql helpers repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for mysql helpers repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for mysql helpers repository")
	}
	if s.Services.Context == nil {
		return errors.New("context service is nil for mysql helpers repository")
	}
	return nil
}
