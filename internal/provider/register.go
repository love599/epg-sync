package provider

import (
	"fmt"
	"sync"

	"github.com/epg-sync/epgsync/internal/model"
	"github.com/epg-sync/epgsync/pkg/errors"
	"github.com/epg-sync/epgsync/pkg/logger"
)

type ProviderFactory func(config *model.ProviderConfig) (Provider, error)

type Registry struct {
	mu        sync.RWMutex
	factories map[string]ProviderFactory
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]ProviderFactory),
	}
}

func (r *Registry) Register(providerType string, factory ProviderFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[providerType]; exists {
		return errors.AlreadyExists("provider_type", providerType)
	}

	r.factories[providerType] = factory

	return nil
}

func (r *Registry) Create(providerType string, config *model.ProviderConfig) (Provider, error) {
	r.mu.RLock()
	factory, exists := r.factories[providerType]
	r.mu.RUnlock()

	if exists {
		logger.Debug("Already registered provider", logger.String("type", providerType))
	}

	if factory == nil {
		return nil, errors.New(
			errors.ErrCodeProviderNotFound,
			fmt.Sprintf("provider type not registered: %s", providerType),
		)
	}

	return factory(config)
}

var globalRegistry = NewRegistry()

func GlobalRegistry() *Registry {
	return globalRegistry
}

func Register(providerType string, factory ProviderFactory) error {
	return globalRegistry.Register(providerType, factory)
}
