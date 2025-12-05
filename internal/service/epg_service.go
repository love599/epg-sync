package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/epg-sync/epgsync/internal/cache"
	"github.com/epg-sync/epgsync/internal/model"
	"github.com/epg-sync/epgsync/internal/provider"
	"github.com/epg-sync/epgsync/internal/repository"
	"github.com/epg-sync/epgsync/pkg/errors"
	"github.com/epg-sync/epgsync/pkg/logger"
)

type EPGService struct {
	programRepo     repository.ProgramRepository
	channelRepo     repository.ChannelRepository
	cache           cache.Cache
	channelMappings repository.ChannelMappingsRepository
	chain           *provider.Chain
}

func NewEPGService(
	programRepo repository.ProgramRepository,
	channelRepo repository.ChannelRepository,
	channelMappings repository.ChannelMappingsRepository,
	cache cache.Cache,
	chain *provider.Chain,
) *EPGService {
	return &EPGService{
		programRepo:     programRepo,
		channelRepo:     channelRepo,
		channelMappings: channelMappings,
		cache:           cache,
		chain:           chain,
	}
}

func (s *EPGService) GetEPG(ctx context.Context, channelID string, date time.Time, page, pageSize int) ([]*model.Program, int64, error) {
	logger.Info("Getting EPG",
		logger.String("channel_id", channelID),
		logger.Time("date", date),
	)

	programs, total, err := s.programRepo.ListByChannelIDAndDate(ctx, channelID, date, page, pageSize)
	if err == nil && len(programs) > 0 {
		return programs, total, nil
	}

	return nil, 0, errors.EPGNotFound(channelID, date.Format("2006-01-02"))
}

func (s *EPGService) GetEPGRange(ctx context.Context, channelID string, startDate, endDate time.Time) ([]*model.Program, error) {
	logger.Info("Getting EPG range",
		logger.String("channel_id", channelID),
		logger.Time("start_date", startDate),
		logger.Time("end_date", endDate),
	)

	programs, err := s.programRepo.ListByChannelIDAndTimeRange(ctx, channelID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return programs, nil
}

func (s *EPGService) GetCurrentProgram(ctx context.Context, channelID string) (*model.Program, error) {
	program, err := s.programRepo.GetCurrentProgram(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func (s *EPGService) SyncEPG(ctx context.Context, channelID string, startDate, endDate time.Time) error {
	logger.Info("Syncing EPG",
		logger.String("channel_id", channelID),
		logger.Time("start_date", startDate),
		logger.Time("end_date", endDate),
	)

	dates := s.getDatesInRange(startDate, endDate)

	for _, date := range dates {
		s.cache.Delete(ctx, "xmltv_epg:"+date.Format("2006-01-02"))
		exists, err := s.programRepo.Exists(ctx, channelID, date)
		if err != nil {
			logger.Error("Failed to check EPG existence",
				logger.Err(err),
				logger.String("channel_id", channelID),
				logger.Time("date", date),
			)
			continue
		}

		if exists {
			logger.Info("EPG already exists, skipping",
				logger.String("channel_id", channelID),
				logger.Time("date", date),
			)
			continue
		}

		channelMaps, err := s.channelMappings.GetByCanonicalID(ctx, channelID)

		logger.Debug("Channel mappings fetched",
			logger.String("channel_id", channelID),
			logger.Int("mapping_count", len(channelMaps)),
		)

		if err != nil {
			logger.Error("Failed to get channel mapping",
				logger.Err(err),
				logger.String("channel_id", channelID),
			)
			return err
		}

		if len(channelMaps) == 0 {
			logger.Warn("No channel mapping found",
				logger.String("channel_id", channelID),
			)
			return errors.ChannelMappingNotFound(channelID)
		}

		var channelMappingInfos []*model.ChannelMappingInfo

		for _, channelMap := range channelMaps {
			channelMappingInfos = append(channelMappingInfos, &model.ChannelMappingInfo{
				ProviderChannelID: channelMap.ProviderChannelID,
				CanonicalID:       channelID,
				ProviderID:        channelMap.ProviderID,
			})
		}
		programs, err := s.chain.FetchEPGParallel(ctx, channelMappingInfos, date)
		if err != nil {
			logger.Warn("Failed to fetch EPG",
				logger.Err(err),
				logger.String("channel_id", channelID),
				logger.Time("date", date),
			)
			continue
		}

		if len(programs) == 0 {
			logger.Warn("No EPG data found",
				logger.String("channel_id", channelID),
				logger.Time("date", date),
			)
			continue
		}

		if err := s.programRepo.CreateBatch(ctx, programs); err != nil {
			logger.Warn("Failed to save EPG",
				logger.Err(err),
				logger.String("channel_id", channelID),
				logger.Time("date", date),
			)
			continue
		}

		cacheKey := s.buildCacheKey(channelID, date)
		s.cache.Delete(ctx, cacheKey)

		logger.Info("Synced EPG",
			logger.String("channel_id", channelID),
			logger.Time("date", date),
			logger.Int("count", len(programs)),
		)
	}

	return nil
}

func (s *EPGService) SyncEPGBatch(ctx context.Context, channelMappingInfos []*model.ChannelMappingInfo, startDate, endDate time.Time, forceUpdate bool) error {
	if len(channelMappingInfos) == 0 {
		logger.Info("No channel mappings provided, skipping EPG sync")
		return nil
	}

	logger.Info("Syncing EPG batch",
		logger.Int("channel_count", len(channelMappingInfos)),
		logger.Time("start_date", startDate),
		logger.Time("end_date", endDate),
		logger.Bool("force_update", forceUpdate),
	)

	dates := s.getDatesInRange(startDate, endDate)
	for _, date := range dates {
		s.cache.Delete(ctx, "xmltv_epg:"+date.Format("2006-01-02"))
		programsToSave := make([]*model.Program, 0)
		cmInfosToSync := make([]*model.ChannelMappingInfo, 0)
		if forceUpdate {
			providerID := channelMappingInfos[0].ProviderID
			logger.Debug("Force update enabled, deleting existing EPG for date",
				logger.Time("date", date),
				logger.String("providerID", providerID),
			)
			if err := s.programRepo.DeleteByDateAndProviderID(ctx, date, providerID); err != nil {
				logger.Warn("Failed to delete existing EPG before update",
					logger.Err(err),
					logger.Time("date", date),
					logger.String("providerID", providerID),
				)
			}
		}

		for _, cmInfo := range channelMappingInfos {
			channelID := cmInfo.CanonicalID
			exists, err := s.programRepo.Exists(ctx, channelID, date)
			if err != nil {
				logger.Warn("Failed to check EPG existence",
					logger.Err(err),
					logger.String("channel_id", channelID),
					logger.Time("date", date),
				)
				continue
			}

			if !exists {
				cmInfosToSync = append(cmInfosToSync, cmInfo)
			} else {
				logger.Debug("EPG already exists, skipping",
					logger.String("channel_id", channelID),
					logger.Time("date", date),
				)
				continue
			}
		}

		if len(cmInfosToSync) == 0 {
			logger.Info("No channels to sync for date",
				logger.Time("date", date),
			)
			continue
		}

		programs, err := s.chain.FetchEPGParallel(ctx, cmInfosToSync, date)

		if err != nil {
			logger.Warn("Failed to fetch EPG batch",
				logger.Err(err),
				logger.Time("date", date),
			)
			continue
		}

		if len(programs) == 0 {
			logger.Warn("No EPG data found",
				logger.Time("date", date),
			)
			continue
		}

		programsToSave = append(programsToSave, programs...)

		if err := s.programRepo.CreateBatch(ctx, programsToSave); err != nil {
			logger.Warn("Failed to save EPG batch",
				logger.Err(err),
				logger.Time("date", date),
				logger.Int("channel_count", len(cmInfosToSync)),
				logger.Int("program_count", len(programsToSave)),
			)
			continue
		}
	}

	logger.Info("Synced EPG batch",
		logger.Int("channel_count", len(channelMappingInfos)),
		logger.Time("start_date", startDate),
		logger.Time("end_date", endDate),
	)

	return nil

}

func (s *EPGService) CleanupOldEPG(ctx context.Context, keepDays int) (int64, error) {
	before := time.Now().AddDate(0, 0, -keepDays)

	logger.Info("Cleaning up old EPG",
		logger.Time("before", before),
		logger.Int("keep_days", keepDays),
	)

	count, err := s.programRepo.DeleteBefore(ctx, before)
	if err != nil {
		return 0, err
	}

	logger.Info("Cleaned up old EPG", logger.Int64("count", count))

	return count, nil
}

func (s *EPGService) buildCacheKey(channelID string, date time.Time) string {
	timezone := date.Location()
	return fmt.Sprintf("epg:%s:%s_%s", channelID, date.Format("2006-01-02"), timezone.String())
}

func (s *EPGService) getDatesInRange(start, end time.Time) []time.Time {
	dates := make([]time.Time, 0)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}

	return dates
}

func (s *EPGService) GenerateXMLTVFormatPrograms(ctx context.Context) (*model.XMLTVEPG, error) {
	logger.Info("Generating XMLTV programs")
	now := time.Now()
	startDate := now.AddDate(0, 0, -1)
	endDate := now.AddDate(0, 0, 1)
	cacheKey := fmt.Sprintf("xmltv_epg:%s", now.Format("2006-01-02"))

	var cachedXMLTV model.XMLTVEPG
	err := s.cache.Get(ctx, cacheKey, &cachedXMLTV)
	if err == nil {
		logger.Info("Returning cached XMLTV EPG",
			logger.Int("channel_count", len(cachedXMLTV.Channels)),
			logger.Int("program_count", len(cachedXMLTV.Programmes)),
		)
		return &cachedXMLTV, nil
	}

	programs, err := s.programRepo.ListAllByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	xmltvEPG := &model.XMLTVEPG{
		Channels:   []*model.XMLTVChannel{},
		Programmes: []*model.XMLTVProgram{},
	}

	channelSet := make(map[string]bool)

	for _, program := range programs {
		if _, exists := channelSet[program.ChannelID]; !exists {
			channelSet[program.ChannelID] = true
			xmltvChannel := &model.XMLTVChannel{
				ID: program.Channel.ChannelID,
				DisplayName: []string{
					program.Channel.ChannelID,
					program.Channel.DisplayName,
				},
			}
			xmltvEPG.Channels = append(xmltvEPG.Channels, xmltvChannel)
		}
		location, err := time.LoadLocation(program.OriginalTimezone)
		if err != nil {
			location = time.UTC
		}
		xmltvProgram := &model.XMLTVProgram{
			Channel: program.Channel.ChannelID,
			Start:   program.StartTime.In(location).Format("20060102150405 -0700"),
			Stop:    program.EndTime.In(location).Format("20060102150405 -0700"),
			Title:   program.Title,
		}
		xmltvEPG.Programmes = append(xmltvEPG.Programmes, xmltvProgram)
	}

	logger.Info("Generated XMLTV programs",
		logger.Int("channel_count", len(xmltvEPG.Channels)),
		logger.Int("program_count", len(xmltvEPG.Programmes)),
	)

	s.cache.Set(ctx, cacheKey, xmltvEPG, 24*time.Hour)

	return xmltvEPG, nil

}

func (s *EPGService) GenerateDIYPFormatPrograms(ctx context.Context, channelName string, date string) (*model.DIYPChannelEPG, error) {
	logger.Info("Generating DIYP format programs")

	channelName = strings.TrimSpace(channelName)

	channel, err := s.channelRepo.GetByChannelName(ctx, channelName)

	if err != nil {
		logger.Error("Failed to get channel by name",
			logger.Err(err),
			logger.String("channel_name", channelName),
		)
		return nil, err
	}

	programs, err := s.programRepo.ListByChannelAndDate(ctx, channel, date)
	if err != nil {
		return nil, err
	}

	var result = &model.DIYPChannelEPG{
		ChannelName: channel.ChannelID,
		Date:        date,
		EPGData:     make([]*model.DIYPProgram, 0),
	}

	for _, p := range programs {
		loc, err := time.LoadLocation(p.OriginalTimezone)
		if err != nil {
			loc = time.UTC
		}
		localStart := p.StartTime.In(loc)
		localEnd := p.EndTime.In(loc)

		result.EPGData = append(result.EPGData, &model.DIYPProgram{
			Start: localStart.Format("15:04"),
			End:   localEnd.Format("15:04"),
			Title: p.Title,
		})
	}

	logger.Info("Generated DIYP programs", logger.Int("count", len(result.EPGData)))
	return result, nil
}
