package middleware

import (
	"backend/pkg/jwt"
	"backend/pkg/redisclient"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware аутентифицирует пользователя по JWT-токену.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var token string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("auth_token")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
		}

		// Проверяем кэш Redis по токену
		cacheKey := "auth_token:" + token
		userID, err := redisclient.Client.Get(redisclient.Ctx, cacheKey).Int()
		if err == nil { // Если нашли в кеше — используем
			c.Set("user_id", userID)
			c.Next()
			return
		}

		// Если в кеше нет, валидируем токен
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Кэшируем user_id с TTL 30 минут
		redisclient.Client.Set(redisclient.Ctx, cacheKey, claims.UserID, 30*time.Minute)

		// Сохраняем user_id в контексте
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
