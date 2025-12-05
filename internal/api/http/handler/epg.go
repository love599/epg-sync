package handler

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/epg-sync/epgsync/internal/api/dto"
	"github.com/epg-sync/epgsync/internal/service"
	"github.com/epg-sync/epgsync/pkg/logger"
	"github.com/gin-gonic/gin"
)

type EPGHandler struct {
	epgService *service.EPGService
}

func NewEPGHandler(epgService *service.EPGService) *EPGHandler {
	return &EPGHandler{
		epgService: epgService,
	}
}
func (h *EPGHandler) GetEPGByChannelAndDate(c *gin.Context) {
	var req dto.GetEPGRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid request parameters", err))
		return
	}
	ctx := c.Request.Context()

	formatDate, err := time.Parse("2006-01-02", req.Date)
	tz := req.Timezone
	if tz == "" {
		tz = "UTC"
	}
	loc, locErr := time.LoadLocation(tz)
	if locErr != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid timezone", locErr))
		return
	}
	formatDate = formatDate.In(loc)
	logger.Debug("Using timezone", logger.String("timezone", tz), logger.Time("date", formatDate))

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid date format", err))
		return
	}

	programs, total, err := h.epgService.GetEPG(ctx, req.ChannelID, formatDate, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to get EPG", err))
		return
	}

	items := make([]any, len(programs))
	for i, p := range programs {
		items[i] = p
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(items, total, req.Page, req.PageSize))
}

func (h *EPGHandler) SyncEPGByChannelAndDateRange(c *gin.Context) {
	var req dto.SyncEPGRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid request parameters", err))
		return
	}
	ctx := c.Request.Context()

	loc, locErr := time.LoadLocation("UTC")
	if locErr != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid timezone", locErr))
		return
	}

	startDate, err := time.ParseInLocation("2006-01-02", req.StartDate, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid start_date format", err))
		return
	}

	endDate, err := time.ParseInLocation("2006-01-02", req.EndDate, loc)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid end_date format", err))
		return
	}

	if err := h.epgService.SyncEPG(ctx, req.ChannelID, startDate, endDate); err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to sync EPG", err))
		return
	}

	c.JSON(http.StatusOK, dto.Success("EPG sync successful"))
}

func (h *EPGHandler) GetEPGByDateRange(c *gin.Context) {
	var req dto.GetEPGByDateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid request parameters", err))
		return
	}
	ctx := c.Request.Context()

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid start_date format", err))
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequest("Invalid end_date format", err))
		return
	}
	programs, err := h.epgService.GetEPGRange(ctx, req.ChannelID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to get EPG range", err))
		return
	}

	c.JSON(http.StatusOK, dto.Success(programs))
}

func (h *EPGHandler) GenerateXMLTVProgram(c *gin.Context) {
	ctx := c.Request.Context()

	xmlData, err := h.epgService.GenerateXMLTVFormatPrograms(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to generate XMLTV", err))
		return
	}

	result, err := xml.MarshalIndent(xmlData, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to marshal XMLTV", err))
		return
	}

	c.Data(http.StatusOK, "application/xml", result)
}

func (h *EPGHandler) GenerateDIYPProgram(c *gin.Context) {

	var req dto.DIYPProgramRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(200, dto.Success(nil))
		return
	}
	ctx := c.Request.Context()

	data, err := h.epgService.GenerateDIYPFormatPrograms(ctx, req.Ch, req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalServerError("Failed to generate DIYP data", err))
		return
	}

	c.JSON(http.StatusOK, data)
}
