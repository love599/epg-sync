package provider

import (
	"fmt"
	"sync"

	"github.com/epg-sync/epgsync/internal/cache"
	"github.com/epg-sync/epgsync/internal/model"
	"github.com/epg-sync/epgsync/pkg/logger"
)

type Factory struct {
	registry  *Registry
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewFactory() *Factory {
	return &Factory{
		registry:  GlobalRegistry(),
		providers: make(map[string]Provider),
	}
}

func (f *Factory) CreateProvider(config *model.ProviderConfig) (Provider, error) {
	f.mu.RLock()
	if p, exists := f.providers[config.ID]; exists {
		f.mu.RUnlock()
		return p, nil
	}
	f.mu.RUnlock()

	provider, err := f.registry.Create(config.ID, config)
	if err != nil {
		return nil, fmt.Errorf("create provider %s: %w", config.ID, err)
	}

	f.mu.Lock()
	f.providers[config.ID] = provider
	f.mu.Unlock()

	return provider, nil
}

func (f *Factory) CreateProviders(configs []model.ProviderConfig) ([]Provider, error) {
	providers := make([]Provider, 0, len(configs))

	for _, cfg := range configs {
		if !cfg.Enabled {
			continue
		}

		p, err := f.CreateProvider(&cfg)
		if err != nil {
			logger.Warn("Failed to create provider", logger.String("id", cfg.ID), logger.Err(err))
			continue
		}

		providers = append(providers, p)
	}

	return providers, nil
}

func (f *Factory) CreateChain(configs []model.ProviderConfig, cache cache.Cache) (*Chain, error) {
	providers, err := f.CreateProviders(configs)
	if err != nil {
		return nil, err
	}

	return NewChain(cache, providers...), nil
}

var globalFactory = NewFactory()

func GlobalFactory() *Factory {
	return globalFactory
}
