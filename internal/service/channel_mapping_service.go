package service

import (
	"context"
	"strings"

	"github.com/epg-sync/epgsync/internal/model"
	"github.com/epg-sync/epgsync/internal/repository"
	"github.com/epg-sync/epgsync/pkg/logger"
)

type ChannelMappingService struct {
	channelRepo        repository.ChannelRepository
	channelMappingRepo repository.ChannelMappingsRepository
}

func NewChannelMappingService(
	channelMappingRepo repository.ChannelMappingsRepository,
	channelRepo repository.ChannelRepository,
) *ChannelMappingService {
	return &ChannelMappingService{
		channelMappingRepo: channelMappingRepo,
		channelRepo:        channelRepo,
	}
}

func (s *ChannelMappingService) AutoMapChannels(ctx context.Context, providerID string, providerChannels []*model.ProviderChannel) error {
	standardChannels, err := s.channelRepo.GetAllChannels(ctx)
	if err != nil {
		return err
	}

	for _, pc := range providerChannels {

		existing, _ := s.channelMappingRepo.GetByProviderChannelID(ctx, pc.ID, providerID)
		if existing != nil && existing.Confidence >= 0.8 {
			continue
		}

		bestMatch, score := s.findBestMatch(pc, standardChannels)
		if bestMatch != nil && score >= 0.8 {
			mapping := &model.ChannelMapping{
				ProviderID:        providerID,
				ProviderChannelID: pc.ID,
				CanonicalID:       bestMatch.ChannelID,
				Confidence:        score,
			}

			if err := s.channelMappingRepo.Create(ctx, mapping); err != nil {
				return err
			}
		} else {
			logger.Info("No suitable channel mapping found",
				logger.String("provider_id", providerID),
				logger.String("provider_channel_id", pc.ID),
				logger.String("provider_channel_name", pc.Name),
				logger.Float64("score", score),
			)
		}
	}

	return nil
}

func (s *ChannelMappingService) ListChannels(ctx context.Context, providerID string) ([]*model.ChannelMapping, error) {
	channels, err := s.channelMappingRepo.ListByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *ChannelMappingService) findBestMatch(pc *model.ProviderChannel, standardChannels []*model.Channel) (*model.Channel, float64) {
	var bestMatch *model.Channel
	var bestScore float64

	for _, sc := range standardChannels {
		score := s.calculateMatchScore(pc, sc)
		if score > bestScore {
			bestScore = score
			bestMatch = sc
		}
	}

	return bestMatch, bestScore
}

func (s *ChannelMappingService) calculateMatchScore(pc *model.ProviderChannel, sc *model.Channel) float64 {
	var score float64

	if strings.EqualFold(pc.ID, sc.ChannelID) || strings.EqualFold(pc.Name, sc.DisplayName) {
		score += 1.0
	}

	if strings.Contains(strings.ToLower(sc.ChannelID), strings.ToLower(pc.Name)) ||
		strings.Contains(strings.ToLower(pc.Name), strings.ToLower(sc.ChannelID)) {
		score += 0.7
	}

	for _, alias := range pc.Aliases {
		if strings.EqualFold(alias, sc.ChannelID) || strings.EqualFold(alias, sc.DisplayName) {
			score += 0.9
			break
		}
	}

	cleanPC := cleanChannelName(pc.Name)
	cleanSC := cleanChannelName(sc.ChannelID)
	cleanSCDisplayName := cleanChannelName(sc.DisplayName)

	if strings.Contains(cleanSCDisplayName, cleanPC) {
		score += 1.0
	}

	if cleanPC == cleanSC {
		score += 0.8
	}

	return score
}

func cleanChannelName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return name
}
