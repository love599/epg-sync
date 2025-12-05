package provider

import (
	"context"
	"time"

	"github.com/epg-sync/epgsync/internal/cache"
	"github.com/epg-sync/epgsync/internal/model"
)

type Provider interface {
	GetID() string

	GetName() string

	FetchEPG(ctx context.Context, providerChannelID, channelID string, date time.Time) ([]*model.Program, error)

	FetchEPGBatch(ctx context.Context, channelMappingInfo []*model.ChannelMappingInfo, date time.Time) ([]*model.Program, error)

	HealthCheck(ctx context.Context) *model.ProviderHealth

	ListChannels() []*model.ProviderChannel

	SupportChannel(providerID, channelID string) bool

	Validate() error

	IsEnabled() bool

	GetPriority() int

	SetCache(cache cache.Cache)

	GetCache() cache.Cache
}
