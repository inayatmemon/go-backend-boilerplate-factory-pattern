package context_repository

import "context"

type Repository interface {
	GetContext() (context.Context, context.CancelFunc)
}

type service struct {
	Input
}

type Input struct{}

func InitService(input Input) Repository {
	return &service{
		input,
	}
}
