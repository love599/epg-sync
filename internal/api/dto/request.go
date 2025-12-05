package dto

type CreateChannelRequest struct {
	ChannelID   string `json:"channel_id" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Category    string `json:"category"`
	Area        string `json:"area"`
	LogoURL     string `json:"logo_url"`
	Regexp      string `json:"regexp"`
	Timezone    string `json:"timezone"`
}

type UpdateChannelRequest struct {
	ID       int64 `json:"id" binding:"required"`
	IsActive int   `json:"is_active"`
	CreateChannelRequest
}

type BatchCreateChannelRequest struct {
	Channels []CreateChannelRequest `json:"channels" binding:"required,dive,required"`
}

type GetEPGRequest struct {
	ChannelID string `form:"channel_id" binding:"omitempty"`
	Date      string `form:"date" binding:"required,datetime=2006-01-02"`
	Timezone  string `form:"timezone" binding:"omitempty"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type GetEPGByDateRangeRequest struct {
	ChannelID string `form:"channel_id" binding:"required"`
	StartDate string `form:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"required,datetime=2006-01-02"`
}

type SyncEPGRequest struct {
	ChannelID string `form:"channel_id" binding:"required"`
	StartDate string `form:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"required,datetime=2006-01-02"`
}

type DIYPProgramRequest struct {
	Ch   string `form:"ch" binding:"required"`
	Date string `form:"date" binding:"required,datetime=2006-01-02"`
}
