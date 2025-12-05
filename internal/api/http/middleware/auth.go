package middleware

import (
	"net/http"
	"strings"

	"github.com/epg-sync/epgsync/pkg/errors"
	"github.com/epg-sync/epgsync/pkg/logger"
	"github.com/epg-sync/epgsync/pkg/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			logger.Warn("Missing JWT token",
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization token",
				"code":  errors.ErrCodeUnauthorized,
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			logger.Warn("Invalid JWT token",
				logger.Err(err),
				logger.String("path", c.Request.URL.Path),
				logger.String("ip", c.ClientIP()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization token",
				"code":  errors.ErrCodeUnauthorized,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
