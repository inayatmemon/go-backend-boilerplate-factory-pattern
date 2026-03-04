package transactions_service

import (
	"errors"
	mongodb_models "go_boilerplate_project/models/databases/mongodb"
	mysql_models "go_boilerplate_project/models/databases/mysql"
	context_repository "go_boilerplate_project/services/context"
	mongodb_helpers_service "go_boilerplate_project/services/helpers/db/mongodb"
	mysql_helpers_service "go_boilerplate_project/services/helpers/db/mysql"
	"log"

	"go.uber.org/zap"
)

type Repository interface {
	RunMongoDBTransaction(input *mongodb_models.TransactionInput) error
	RunMySQLTransaction(input *mysql_models.TransactionInput) error
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
	MongoDB mongodb_helpers_service.Repository
	MySQL   mysql_helpers_service.Repository
}

type Services struct {
	Context context_repository.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for transactions repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for transactions repository")
	}
	if s.Helpers == nil {
		return errors.New("helpers is nil for transactions repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for transactions repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for transactions repository")
	}
	if s.Helpers.MongoDB == nil {
		return errors.New("mongodb helpers is nil for transactions repository")
	}
	if s.Helpers.MySQL == nil {
		return errors.New("mysql helpers is nil for transactions repository")
	}
	if s.Services.Context == nil {
		return errors.New("context is nil for transactions repository")
	}
	return nil
}
