package middleware

import (
	"backend/pkg/i18n"
	"strings"

	"github.com/gin-gonic/gin"
)

// I18nMiddleware извлекает язык из заголовка Accept-Language или X-Language и помещает его в контекст.
// Поддерживаемые языки: "en", "ru". Если язык не поддерживается или отсутствует, используется "ru".
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("X-Service-Language")
		lang = strings.ToLower(lang)
		if lang != "en" && lang != "ru" {
			lang = "ru" // Default language
		}
		// Create a new RequestLanguageProvider for this request
		provider := &i18n.RequestLanguageProvider{Lang: lang} // Use exported field Lang
		i18n.SetLanguageProvider(provider)
		// Optionally keep it in Gin context for handlers
		c.Set("lang", lang)
		c.Next()
	}
}
