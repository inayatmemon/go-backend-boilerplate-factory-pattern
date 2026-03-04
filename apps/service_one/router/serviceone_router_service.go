package serviceone_router

import (
	"fmt"
	api_constants "go_boilerplate_project/constants/api"
	"os"

	"github.com/gin-gonic/gin"
)

func (s *service) ConfigureRouter() {
	s.Engine.Engine.Use(gin.Recovery())
	s.Engine.Engine.Use(gin.Logger())

	// Apply global middlewares to all routes
	for _, mw := range s.Input.Middleware.Global.GetMiddlewares() {
		s.Engine.Engine.Use(mw)
	}
}

func (s *service) SetupRoutes() {
	// v1 group with application-specific middlewares
	v1 := s.Engine.Engine.Group(api_constants.BaseURLV1)
	for _, mw := range s.Input.Middleware.Application.GetMiddlewares() {
		v1.Use(mw)
	}
	{
		v1.POST("/brands", s.Http.Brands.EditBrand)
		v1.POST("/products", s.Http.Products.CreateProduct)
	}
}

func (s *service) Run() {
	err := s.Engine.Engine.Run(fmt.Sprintf(":%d", s.Input.Env.App.AppPort))
	if err != nil {
		s.Input.Logger.Error("Failed to run serviceone router", "error", err)
		os.Exit(1)
	}
}
