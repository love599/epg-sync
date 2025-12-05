package handler

import (
	"net/http"

	"github.com/epg-sync/epgsync/internal/api/dto"
	"github.com/epg-sync/epgsync/internal/model"
	"github.com/epg-sync/epgsync/internal/service"
	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	channelService *service.ChannelService
}

func NewChannelHandler(channelService *service.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
	}
}

func (h *ChannelHandler) CreateChannel(c *gin.Context) {
	var channel *dto.CreateChannelRequest
	if err := c.ShouldBindBodyWithJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid request parameters", err))
		return
	}
	ctx := c.Request.Context()

	newChannel := &model.Channel{
		ChannelID:   channel.ChannelID,
		DisplayName: channel.DisplayName,
		Category:    channel.Category,
		Area:        channel.Area,
		Regexp:      channel.Regexp,
		LogoURL:     channel.LogoURL,
		Timezone:    channel.Timezone,
	}

	createdChannel, err := h.channelService.CreateChannel(ctx, newChannel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to create channel", err))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(createdChannel))

}

func (h *ChannelHandler) BatchCreateChannel(c *gin.Context) {
	var channels dto.BatchCreateChannelRequest
	if err := c.ShouldBindBodyWithJSON(&channels); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid request parameters", err))
		return
	}
	ctx := c.Request.Context()

	var newChannels []*model.Channel
	for _, channel := range channels.Channels {
		newChannels = append(newChannels, &model.Channel{
			ChannelID:   channel.ChannelID,
			DisplayName: channel.DisplayName,
			Category:    channel.Category,
			Area:        channel.Area,
			LogoURL:     channel.LogoURL,
			Timezone:    channel.Timezone,
		})
	}

	createdChannels, err := h.channelService.BatchCreateChannels(ctx, newChannels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to create channels", err))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(createdChannels))
}

func (h *ChannelHandler) ListChannels(c *gin.Context) {
	ctx := c.Request.Context()

	channels, err := h.channelService.ListAllChannels(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to list channels", err))
		return
	}

	c.JSON(http.StatusOK, dto.Success(channels))
}

func (h *ChannelHandler) GetChannel(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	channel, err := h.channelService.GetChannel(ctx, channelID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.NotFound("Channel not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(channel))
}

func (h *ChannelHandler) UpdateChannel(c *gin.Context) {
	channelID := c.Param("id")
	var req dto.UpdateChannelRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid request parameters", err))
		return
	}

	ctx := c.Request.Context()
	channel := &model.Channel{
		ID:          req.ID,
		ChannelID:   channelID,
		DisplayName: req.DisplayName,
		Category:    req.Category,
		IsActive:    req.IsActive,
		Area:        req.Area,
		LogoURL:     req.LogoURL,
		Regexp:      req.Regexp,
		Timezone:    req.Timezone,
	}

	if err := h.channelService.UpdateChannel(ctx, channel); err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to update channel", err))
		return
	}

	c.JSON(http.StatusOK, dto.Success(channel))
}

func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	if err := h.channelService.DeleteChannel(ctx, channelID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to delete channel", err))
		return
	}

	c.JSON(http.StatusOK, dto.Success(gin.H{"message": "Channel deleted successfully"}))
}

func (h *ChannelHandler) ListChannelMappings(c *gin.Context) {
	ctx := c.Request.Context()

	mappings, err := h.channelService.ListAllChannelMappings(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to list channel mappings", err))
		return
	}

	c.JSON(http.StatusOK, dto.Success(mappings))
}

func (h *ChannelHandler) GetChannelMappings(c *gin.Context) {
	channelID := c.Param("id")
	ctx := c.Request.Context()

	mappings, err := h.channelService.GetChannelMappings(ctx, channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to get channel mappings", err))
		return
	}

	c.JSON(http.StatusOK, dto.Success(mappings))
}
