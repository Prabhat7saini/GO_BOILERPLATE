package middleware

import (
	"boiler-platecode/src/common/const/exception"
	"boiler-platecode/src/common/utils"
	"boiler-platecode/src/core/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
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

		log.Printf("token decodevalue %v\n", claims)
		
		// Store user ID in context
		if userID, ok := claims["userId"].(float64); ok {
			c.Set("userID", int(userID))
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}


		c.Next()
	}
}
