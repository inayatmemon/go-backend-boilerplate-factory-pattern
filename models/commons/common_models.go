package common_models

import "context"

type IsContextPresentInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
}
