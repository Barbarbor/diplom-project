package middleware

import (
	"backend/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Попробуем получить токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")

		var token string
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Если в заголовке токена нет, пробуем получить из cookie
			var err error
			token, err = c.Cookie("auth_token")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
		}
		// Проверяем токен
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// Сохраняем данные пользователя в контексте запроса
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
