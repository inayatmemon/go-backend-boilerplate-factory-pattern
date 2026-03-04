package application_middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	HeaderAppName    = "X-App-Name"
	HeaderAppVersion = "X-App-Version"
)

// AppVersion adds application name and version to response headers for service-specific routes.
func (s *service) AppVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.Input.AppName != "" {
			c.Header(HeaderAppName, s.Input.AppName)
		}
		if s.Input.AppVersion != "" {
			c.Header(HeaderAppVersion, s.Input.AppVersion)
		}
		c.Next()
	}
}

// GetMiddlewares returns all application-specific middlewares.
func (s *service) GetMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		s.AppVersion(),
	}
}
