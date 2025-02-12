package access

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"

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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "access denied"})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !strings.HasPrefix(authHeader, authPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		accessToken := strings.TrimPrefix(authHeader, authPrefix)

		claims, err := jwt.VerifyToken(accessToken, m.authConfig.AccessTokenSecretKey())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		user, err := m.userService.GetUserByName(c, claims.Username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
