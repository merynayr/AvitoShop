package access

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"
)

// Middleware структура для проверки доступа
type Middleware struct {
	accessService service.AccessService
	authConfig    config.AuthConfig
}

// NewAccessMiddleware возвращает новый объект middleware слоя access
func NewAccessMiddleware(accessService service.AccessService, authConfig config.AuthConfig) *Middleware {
	return &Middleware{
		accessService: accessService,
		authConfig:    authConfig,
	}
}

// Check проверяет доступ к ресурсу
func (m *Middleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()

		user, err := m.accessService.Check(c, endpoint)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// AddAccessTokenFromCookie извлекает токен из cookie и добавляет в заголовок Authorization
func (m *Middleware) AddAccessTokenFromCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			accessToken, err := c.Cookie("access_token")
			if err == nil {
				c.Request.Header.Set("Authorization", "Bearer "+accessToken)
			}
		}
		c.Next()
	}
}
