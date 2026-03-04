package products_http_service

import (
	products_models "go_boilerplate_project/models/products"

	"github.com/gin-gonic/gin"
)

func (s *service) CreateProduct(c *gin.Context) {
	s.Input.Logger.Infow("CreateProduct HTTP request received")

	request := &products_models.CreateProductRequest{}
	if err := s.Input.Helpers.API.ParseJSONBody(c, request); err != nil {
		s.Input.Logger.Warnw("CreateProduct validation failed", "errors", err)
		s.Input.Helpers.API.SendValidationError(c, err)
		return
	}

	response := s.Input.Domain.Products.CreateProduct(request)
	s.Input.Helpers.API.SendApiResponse(c, response)
}
