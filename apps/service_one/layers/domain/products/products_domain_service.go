package products_domain_service

import (
	api_models "go_boilerplate_project/models/api"
	brands_models "go_boilerplate_project/models/brands"
	products_models "go_boilerplate_project/models/products"
	"net/http"
)

func (s *service) CreateProduct(request *products_models.CreateProductRequest) *api_models.ApiResponse {
	s.Input.Logger.Infow("CreateProduct started", "name", request.Name, "brandId", request.BrandId)

	response := &api_models.ApiResponse{
		StatusCode: http.StatusCreated,
		Message:    "Product created successfully",
		Data:       nil,
	}

	// create product in msql
	product := &products_models.Product{
		Name:    request.Name,
		BrandId: request.BrandId,
	}

	brand := &brands_models.Brand{ID: request.BrandId}
	// get brand name from brand id
	err := s.Input.Data.Brands.GetBrandMySQL(brand)
	if err != nil {
		s.Input.Logger.Errorw("CreateProduct failed to get brand", "error", err, "brandId", request.BrandId)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "failed to create product"
		response.Error = err.Error()
		return response
	}

	product.BrandName = brand.Name

	if err := s.Input.Data.Products.CreateProductMySQL(product); err != nil {
		s.Input.Logger.Errorw("CreateProduct failed at MySQL create", "error", err, "name", request.Name)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "failed to create product"
		response.Error = err.Error()
		return response
	}

	// create product in mongodb
	if err := s.Input.Data.Products.CreateProductMongoDB(product); err != nil {
		s.Input.Logger.Errorw("CreateProduct failed at MongoDB create", "error", err, "name", request.Name)
		response.StatusCode = http.StatusInternalServerError
		response.Message = "failed to create product"
		response.Error = err.Error()
		return response
	}

	s.Input.Logger.Infow("CreateProduct completed successfully", "name", request.Name, "brandId", request.BrandId)
	return response
}
