package app

import (
	"context"

	"github.com/epg-sync/epgsync/pkg/logger"
)

func (a *App) autoMapProviderChannels() error {
	ctx := context.Background()

	providers := a.providerChain.GetProviders()

	for _, provider := range providers {

		channels := provider.ListChannels()
		if err := a.services.ChannelMapping.AutoMapChannels(ctx, provider.GetID(), channels); err != nil {
			logger.Warn("Failed to auto map channels",
				logger.String("provider", provider.GetID()),
				logger.Err(err))
		} else {
			logger.Debug("Auto mapped channels",
				logger.String("provider", provider.GetID()),
				logger.Int("count", len(channels)))
		}
	}

	return nil
}
