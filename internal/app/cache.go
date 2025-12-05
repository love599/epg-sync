package app

import (
	"context"
	"time"

	"github.com/epg-sync/epgsync/internal/cache"
	"github.com/epg-sync/epgsync/pkg/logger"
)

func (app *App) initializeCache() error {
	cacheCfg := app.cfg.Cache
	logger.Debug("Initializing redis...",
		logger.String("type", cacheCfg.Type),
		logger.String("address", cacheCfg.Addr),
	)

	switch cacheCfg.Type {
	case "memory":
		app.cache = cache.NewMemoryCache()
	case "redis":
		app.cache = cache.NewRedisCache(cacheCfg.Addr, cacheCfg.Password, cacheCfg.DB)
	default:
		logger.Warn("Unknown cache type, falling back to memory", logger.String("type", cacheCfg.Type))
		app.cache = cache.NewMemoryCache()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.cache.Ping(ctx); err != nil {
		return err
	}

	return nil
}
