package handler

import (
	"net/http"

	"github.com/epg-sync/epgsync/internal/api/dto"
	"github.com/epg-sync/epgsync/internal/service"
	"github.com/gin-gonic/gin"
)

type SchedulerHandler struct {
	schedulerService *service.SchedulerService
}

func NewSchedulerHandler(schedulerService *service.SchedulerService) *SchedulerHandler {
	return &SchedulerHandler{
		schedulerService: schedulerService,
	}
}

func (h *SchedulerHandler) SyncAllEPG(c *gin.Context) {
	forceUpdate := c.Query("force") == "true"

	go h.schedulerService.SyncAllEPG(forceUpdate)

	c.JSON(http.StatusOK, dto.Success("EPG sync job started"))
}
