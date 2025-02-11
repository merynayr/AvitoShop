package access

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/service"
)

// Middleware access структура
type Middleware struct {
	accessService service.AccessService
}

// NewMiddleware возвращает новый объект middleware слоя access
func NewMiddleware(accessService service.AccessService) *Middleware {
	return &Middleware{
		accessService: accessService,
	}
}

// Check - общий middleware для проверки доступа
func (m *Middleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.Request.URL.Path

		username, err := m.accessService.Check(c, endpoint)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			c.Abort()
			return
		}

		c.Set("username", username)

		c.Next()
	}
}
