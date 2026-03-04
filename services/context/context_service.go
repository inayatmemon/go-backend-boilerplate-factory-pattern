package context_repository

import (
	"context"
	numeric_constants "go_boilerplate_project/constants/numeric"
	"time"
)

func (s *service) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), numeric_constants.DefaultContextTimeoutInSeconds*time.Second)
}
