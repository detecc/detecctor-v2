package middleware

import (
	"context"
	"fmt"
	"log"
	"sync"
)

var middlewareManager *Manager

func init() {
	once := sync.Once{}
	once.Do(func() {
		GetMiddlewareManager()
	})
}

// GetMiddlewareManager return a global middleware manager struct (singleton).
func GetMiddlewareManager() *Manager {
	if middlewareManager == nil {
		middlewareManager = &Manager{middlewareMap: sync.Map{}}
	}
	return middlewareManager
}

// Register the middleware plugin in the manager.
func (m *Manager) Register(name string, action Handler) {
	log.Println("Adding middleware", name, action)
	m.middlewareMap.Store(name, action)
}

// HasMiddleware Check if the manager has the middleware stored.
func (m *Manager) HasMiddleware(name string) bool {
	_, isFound := m.middlewareMap.Load(name)
	if !isFound {
		return false
	}
	return true
}

// GetMiddleware gets a middleware stored in the map.
func (m *Manager) GetMiddleware(name string) (Handler, error) {
	middleware, isFound := m.middlewareMap.Load(name)
	if !isFound {
		return nil, fmt.Errorf("middleware %s not found", name)
	}
	return middleware.(Handler), nil
}

// GetAllMiddleware gets all the middleware stored in the manager.
func (m *Manager) GetAllMiddleware() []Handler {
	var middlewares []Handler
	m.middlewareMap.Range(func(key, value interface{}) bool {
		middlewares = append(middlewares, value.(Handler))
		return true
	})
	return middlewares
}

// Chain the middleware in consecutive order. This is useful for processing requests depending on the business constraints.
// Returns an error if it occurred during execution.
func (m *Manager) Chain(ctx context.Context, middleware ...string) error {
	var finalMiddleware Handler
	if len(middleware) == 0 {
		return nil
	}

	for i, key := range middleware {
		if !m.HasMiddleware(key) {
			return fmt.Errorf("middleware %s not found", key)
		}

		m, _ := m.middlewareMap.Load(key)
		if i == 0 {
			// assign the first instance
			finalMiddleware = m.(Handler)
		} else {
			// try to chain the middleware, stop execution if an error occurs
			mw, err := finalMiddleware.Chain(ctx, m.(Handler))
			if err != nil {
				return err
			}
			finalMiddleware = mw
		}
	}

	return finalMiddleware.Execute(ctx)
}
