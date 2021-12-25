package middleware

import (
	"context"
	"sync"
)

type (
	// Handler is the interface that needs to be implemented in order to make functioning middleware
	Handler interface {
		// Execute should be called in the last
		Execute(ctx context.Context) error
		// Chain does some logic and returns the next middleware in the chain and an error, if one occurs during execution
		Chain(ctx context.Context, next Handler) (Handler, error)
	}

	// Manager registers, stores and chains middleware.
	Manager struct {
		middlewareMap sync.Map
	}
)
