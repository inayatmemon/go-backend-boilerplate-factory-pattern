package brands_http_service

import (
	brands_models "go_boilerplate_project/models/brands"

	"github.com/gin-gonic/gin"
)

func (s *service) EditBrand(c *gin.Context) {
	s.Input.Logger.Infow("EditBrand HTTP request received")

	request := &brands_models.EditBrandRequest{}
	if err := s.Helpers.API.ParseJSONBody(c, request); err != nil {
		s.Input.Logger.Warnw("EditBrand validation failed", "errors", err)
		s.Helpers.API.SendValidationError(c, err)
		return
	}

	response := s.Input.Domain.Brands.EditBrand(request.Name, request.ID)
	s.Input.Helpers.API.SendApiResponse(c, response)
}
