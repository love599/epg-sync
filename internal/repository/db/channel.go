// pkg/repository/mysql/channel.go
package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/epg-sync/epgsync/internal/model"
	"github.com/epg-sync/epgsync/internal/repository"
	"github.com/epg-sync/epgsync/pkg/errors"
	"github.com/epg-sync/epgsync/pkg/logger"
	"gorm.io/gorm"
)

type channelRepo struct {
	*BaseRepository
}

func NewChannelRepository(db *gorm.DB) repository.ChannelRepository {
	return &channelRepo{BaseRepository: NewBaseRepository(db)}
}

func (r *channelRepo) Create(ctx context.Context, channel *model.Channel) error {
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(channel).Error; err != nil {
		logger.Error("Failed to create channel",
			logger.Err(err),
			logger.String("channel_id", channel.ChannelID),
		)
		return errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to create channel")
	}

	return nil
}

func (r *channelRepo) CreateBatch(ctx context.Context, channels []*model.Channel) error {
	if len(channels) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(channels).Error; err != nil {
		logger.Error("Failed to create channels in batch",
			logger.Err(err),
			logger.Int("count", len(channels)),
		)
		return errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to create channels in batch")
	}

	logger.Info("Batch created channels", logger.Int("count", len(channels)))
	return nil
}

func (r *channelRepo) GetAllChannels(ctx context.Context) ([]*model.Channel, error) {
	var channels []*model.Channel
	err := r.db.WithContext(ctx).Find(&channels).Error
	if err != nil {
		logger.Error("Failed to get all channels", logger.Err(err))
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to get all channels")
	}

	return channels, nil
}

func (r *channelRepo) GetByID(ctx context.Context, channelID string) (*model.Channel, error) {
	var channel model.Channel
	err := r.db.WithContext(ctx).Where("channel_id = ?", channelID).First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ChannelNotFound(channelID)
		}
		logger.Error("Failed to get channel",
			logger.Err(err),
			logger.String("channel_id", channelID),
		)
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to get channel")
	}

	return &channel, nil
}

func (r *channelRepo) ListByProviderID(ctx context.Context, providerID string) ([]*model.Channel, error) {
	var channels []*model.Channel
	err := r.db.WithContext(ctx).
		Where("provider_id = ?", providerID).
		Order("number ASC").
		Find(&channels).Error
	if err != nil {
		logger.Error("Failed to list channels by provider",
			logger.Err(err),
			logger.String("provider_id", providerID),
		)
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to list channels")
	}

	return channels, nil
}

func (r *channelRepo) GetByChannelName(ctx context.Context, channelName string) (*model.Channel, error) {
	var channel model.Channel
	err := r.db.WithContext(ctx).Where("? REGEXP `regexp`", channelName).First(&channel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ChannelNotFound(channelName)
		}
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to get channel by name")
	}

	return &channel, nil

}

func (r *channelRepo) Search(ctx context.Context, query string, opts *repository.ListOptions) ([]*model.Channel, error) {
	if opts == nil {
		opts = &repository.ListOptions{
			Page:     1,
			PageSize: 50,
			OrderBy:  "name",
			Order:    "asc",
		}
	}

	searchPattern := "%" + query + "%"
	offset := (opts.Page - 1) * opts.PageSize

	var channels []*model.Channel
	err := r.db.WithContext(ctx).
		Where("name LIKE ? OR `group` LIKE ?", searchPattern, searchPattern).
		Order(fmt.Sprintf("%s %s", opts.OrderBy, opts.Order)).
		Limit(opts.PageSize).
		Offset(offset).
		Find(&channels).Error
	if err != nil {
		logger.Error("Failed to search channels",
			logger.Err(err),
			logger.String("query", query),
		)
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to search channels")
	}

	return channels, nil
}

func (r *channelRepo) Update(ctx context.Context, channel *model.Channel) error {
	channel.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).
		Model(&model.Channel{}).
		Where("id = ?", channel.ID).
		Updates(map[string]any{
			"channel_id":   channel.ChannelID,
			"display_name": channel.DisplayName,
			"logo_url":     channel.LogoURL,
			"category":     channel.Category,
			"area":         channel.Area,
			"regexp":       channel.Regexp,
			"is_active":    channel.IsActive,
			"timezone":     channel.Timezone,
			"updated_at":   channel.UpdatedAt,
		})
	if result.Error != nil {
		logger.Error("Failed to update channel",
			logger.Err(result.Error),
			logger.String("channel_id", channel.ChannelID),
		)
		return errors.Wrap(result.Error, errors.ErrCodeDatabaseQuery, "failed to update channel")
	}

	if result.RowsAffected == 0 {
		return errors.ChannelNotFound(channel.ChannelID)
	}

	return nil
}

func (r *channelRepo) Delete(ctx context.Context, channelID string) error {
	result := r.db.WithContext(ctx).Where("channel_id = ?", channelID).Delete(&model.Channel{})
	if result.Error != nil {
		logger.Error("Failed to delete channel",
			logger.Err(result.Error),
			logger.String("channel_id", channelID),
		)
		return errors.Wrap(result.Error, errors.ErrCodeDatabaseQuery, "failed to delete channel")
	}

	if result.RowsAffected == 0 {
		return errors.ChannelNotFound(channelID)
	}

	return nil
}

func (r *channelRepo) DeleteByProviderID(ctx context.Context, providerID string) error {
	result := r.db.WithContext(ctx).
		Where("provider_id = ?", providerID).
		Delete(&model.Channel{})
	if result.Error != nil {
		logger.Error("Failed to delete channels by provider",
			logger.Err(result.Error),
			logger.String("provider_id", providerID),
		)
		return errors.Wrap(result.Error, errors.ErrCodeDatabaseQuery, "failed to delete channels")
	}

	logger.Info("Deleted channels by provider",
		logger.String("provider_id", providerID),
		logger.Int64("count", result.RowsAffected),
	)

	return nil
}

func (r *channelRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Channel{}).Count(&count).Error
	if err != nil {
		logger.Error("Failed to count channels", logger.Err(err))
		return 0, errors.Wrap(err, errors.ErrCodeDatabaseQuery, "failed to count channels")
	}

	return count, nil
}
