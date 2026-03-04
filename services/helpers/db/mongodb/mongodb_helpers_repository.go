package mongodb_helpers_service

import (
	"errors"
	mongodb_models "go_boilerplate_project/models/databases/mongodb"
	context_repository "go_boilerplate_project/services/context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	InsertOne(input *mongodb_models.InsertOneInput) (*mongodb_models.InsertOneOutput, error)
	InsertMany(input *mongodb_models.InsertManyInput) (*mongodb_models.InsertManyOutput, error)

	UpdateOne(input *mongodb_models.UpdateOneInput) (*mongodb_models.UpdateOutput, error)
	UpdateMany(input *mongodb_models.UpdateManyInput) (*mongodb_models.UpdateOutput, error)
	UpdateOneWithID(input *mongodb_models.UpdateOneWithIDInput) (*mongodb_models.UpdateOutput, error)

	DeleteOne(input *mongodb_models.DeleteOneInput) (*mongodb_models.DeleteOneOutput, error)
	DeleteMany(input *mongodb_models.DeleteManyInput) (*mongodb_models.DeleteManyOutput, error)

	FindOne(input *mongodb_models.FindOneInput) error
	FindMany(input *mongodb_models.FindManyInput) error
	FindManyWithFilters(input *mongodb_models.FindManyWithFiltersInput) error

	Aggregate(input *mongodb_models.AggregateInput) error

	RunTransaction(input *mongodb_models.TransactionInput) error
}

type service struct {
	Input Input
}

type Input struct {
	Client   *Client
	Logger   *zap.SugaredLogger
	Services *Services
}

type Services struct {
	Context context_repository.Repository
}

type Client struct {
	MongoDBClient *MongoClient
}

type MongoClient struct {
	Client   *mongo.Client
	Database string
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for mongodb helpers repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for mongodb helpers repository")
	}
	if s.Client == nil {
		return errors.New("client is nil for mongodb helpers repository")
	}
	if s.Client.MongoDBClient == nil {
		return errors.New("mongodb client is nil for mongodb helpers repository")
	}
	if s.Client.MongoDBClient.Client == nil {
		return errors.New("mongodb client is nil for mongodb helpers repository")
	}
	if s.Client.MongoDBClient.Database == "" {
		return errors.New("database is required for mongodb helpers repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for mongodb helpers repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for mongodb helpers repository")
	}
	if s.Services.Context == nil {
		return errors.New("context is nil for mongodb helpers repository")
	}
	return nil
}
