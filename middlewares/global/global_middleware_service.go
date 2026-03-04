package global_middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "X-Request-ID"

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// RequestID adds a unique request ID to each incoming request for tracing and debugging.
// The ID is set in the request context and as a response header.
func (s *service) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDKey, requestID)
		c.Next()
	}
}

// GetMiddlewares returns all global middlewares applied to every service route.
func (s *service) GetMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		s.RequestID(),
	}
}
