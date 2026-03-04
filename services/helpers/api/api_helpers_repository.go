package api_helpers_service

import (
	"errors"
	api_models "go_boilerplate_project/models/api"
	multilang_service "go_boilerplate_project/services/multilang"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Repository interface {
	ParseJSONBody(c *gin.Context, dest any) []api_models.ValidationError
	ParseQueryParams(c *gin.Context, dest any) []api_models.ValidationError
	ParsePathParam(c *gin.Context, key string) string

	SendSuccess(c *gin.Context, statusCode int, message string, data any)
	SendError(c *gin.Context, statusCode int, message string, errDetail string)
	SendValidationError(c *gin.Context, validationErrors []api_models.ValidationError)
	SendPaginatedSuccess(c *gin.Context, statusCode int, message string, data any, pagination *api_models.PaginationResult)
	SendApiResponse(c *gin.Context, response *api_models.ApiResponse)

	CalculatePagination(input *api_models.PaginationInput, totalCount int) *api_models.PaginationResult
	GetPaginationFromQuery(c *gin.Context) (*api_models.PaginationInput, []api_models.ValidationError)
}

type service struct {
	Input Input
}

type Input struct {
	Logger   *zap.SugaredLogger
	Services *Services
}

type Services struct {
	Multilang multilang_service.Repository
}

func InitService(input Input) Repository {
	if err := input.validateInput(); err != nil {
		log.Fatalf("Failed to validate input for api helpers repository: %v", err)
	}
	return &service{
		Input: input,
	}
}

func (s *Input) validateInput() error {
	if s == nil {
		return errors.New("input is nil for api helpers repository")
	}
	if s.Logger == nil {
		return errors.New("logger is nil for api helpers repository")
	}
	if s.Services == nil {
		return errors.New("services is nil for api helpers repository")
	}
	if s.Services.Multilang == nil {
		return errors.New("multilang is nil for api helpers repository")
	}
	return nil
}
