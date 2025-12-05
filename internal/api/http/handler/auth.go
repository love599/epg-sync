package handler

import (
	"net/http"

	"github.com/epg-sync/epgsync/internal/service"
	"github.com/epg-sync/epgsync/pkg/errors"
	"github.com/epg-sync/epgsync/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6,max=100"`
		Email    string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  errors.ErrCodeInvalidParam,
		})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		logger.Error("Failed to register user",
			logger.Err(err),
			logger.String("username", req.Username),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  errors.GetCode(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  errors.ErrCodeInvalidParam,
		})
		return
	}

	token, user, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		logger.Warn("Login failed",
			logger.String("username", req.Username),
			logger.String("ip", c.ClientIP()),
			logger.Err(err),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
			"code":  errors.GetCode(err),
		})
		return
	}

	logger.Info("User logged in",
		logger.String("username", user.Username),
		logger.String("ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
			"code":  errors.ErrCodeUnauthorized,
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
			"code":  errors.ErrCodeNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  errors.ErrCodeInvalidParam,
		})
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID.(int64), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  errors.GetCode(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password changed successfully",
	})
}
