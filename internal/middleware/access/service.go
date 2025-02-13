package access

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"

	"github.com/merynayr/AvitoShop/internal/utils/jwt"
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
		endpoint := c.Request.URL.Path

		username, err := m.accessService.Check(c, endpoint)
		if err != nil {
			sys.HandleError(c, sys.NewCommonError("access denied", codes.Unauthorized))
			c.Abort()
			return
		}

		c.Set("username", username)

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

const (
	authHeader = "Authorization"
	authPrefix = "Bearer "
)

// ExtractUserID получает имя из токена
func (m *Middleware) ExtractUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sys.HandleError(c, sys.NewCommonError("invalid token", codes.Unauthorized))
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, authPrefix) {
			sys.HandleError(c, sys.NewCommonError("invalid token", codes.Unauthorized))
			c.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, authPrefix)

		claims, err := jwt.VerifyToken(accessToken, m.authConfig.AccessTokenSecretKey())
		if err != nil {
			sys.HandleError(c, sys.NewCommonError("invalid token", codes.Unauthorized))
			c.Abort()
			return
		}

		user, err := m.userService.GetUserByName(c, claims.Username)
		if err != nil {
			sys.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
