package middlewares

import (
	"boiler-platecode/src/common/const/exception"
	"boiler-platecode/src/common/lib/logger"
	"boiler-platecode/src/common/utils"
	"boiler-platecode/src/core/config"
	"boiler-platecode/src/core/redis"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	module = "Auth Middleware"
)

func AuthMiddleware(redisService redis.RedisService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
			logger.Warning(module, "AuthMiddlerWare", "Token not found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": exception.GetException(exception.INTERNAL_SERVER_ERROR).Message,
			})
			return
		}

		// Verify token
		claims, err := utils.ValidateJwtToken(token, config.AppConfig.AuthJwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": exception.GetException(exception.INTERNAL_SERVER_ERROR).Message,
			})
			return
		}

		// Store user ID in context
		if userIDFloat, ok := claims["userId"].(float64); ok {
			userID := int(userIDFloat)
			storeToken, err := redisService.Get(fmt.Sprintf("Auth:userId:%d", int(userID)))
			if err != nil || storeToken != token {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized access",
				})
				return
			}
			c.Set("userID", int(userID))
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
			})
			return
		}

		c.Next()
	}
}
