package access

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"
	"github.com/merynayr/AvitoShop/internal/sys"
)

// Middleware access структура
type Middleware struct {
	accessService service.AccessService
	authConfig    config.AuthConfig
}

// NewMiddleware возвращает новый объект middleware слоя access
func NewMiddleware(accessService service.AccessService, authConfig config.AuthConfig) *Middleware {
	return &Middleware{
		accessService: accessService,
		authConfig:    authConfig,
	}
}

// Check - общий middleware для проверки доступа
func (m *Middleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()

		user, err := m.accessService.Check(c, endpoint)
		if err != nil {
			sys.HandleError(c, sys.AccessDeniedError)
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

// AddAccessTokenFromCookie - middleware, который извлекает токен доступа из куки и добавляет его в заголовок
func (m *Middleware) AddAccessTokenFromCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.Next()
			return
		}

		c.Request.Header.Set("Authorization", "Bearer "+accessToken)
		c.Next()
	}
}
