package access

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

// Middleware access структура
type Middleware struct {
	userService   service.UserService
	accessService service.AccessService
	authConfig    config.AuthConfig
}

// NewMiddleware возвращает новый объект middleware слоя access
func NewMiddleware(userService service.UserService, accessService service.AccessService, authConfig config.AuthConfig) *Middleware {
	return &Middleware{
		userService:   userService,
		accessService: accessService,
		authConfig:    authConfig,
	}
}

// Check - общий middleware для проверки доступа
func (m *Middleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()

		username, err := m.accessService.Check(c, endpoint)
		if err != nil {
			sys.HandleError(c, sys.NewCommonError("access denied", codes.Unauthorized))
			c.Abort()
			return
		}

		if username != "" {
			user, err := m.userService.GetUserByName(c, username)

			if err != nil {
				sys.HandleError(c, err)
				c.Abort()
				return
			}

			c.Set("user", user)
		}

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
