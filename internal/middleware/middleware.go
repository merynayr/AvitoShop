package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/middleware/access"
	"github.com/merynayr/AvitoShop/internal/middleware/metrics"
	"github.com/merynayr/AvitoShop/internal/service"
)

// Middleware интерфейс для всех middleware
type Middleware interface {
	Access() *access.Middleware
	Metrics() *metrics.Middleware
	TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc
}

// provider структура, реализующая Middleware
type provider struct {
	accessMiddleware  *access.Middleware
	metricsMiddleware *metrics.Middleware
}

// NewMiddlewareProvider создает новый экземпляр провайдера middleware
func NewMiddlewareProvider(accessService service.AccessService, authConfig config.AuthConfig) Middleware {
	return &provider{
		accessMiddleware:  access.NewAccessMiddleware(accessService, authConfig),
		metricsMiddleware: metrics.NewMetricsMiddleware(),
	}
}

// Access возвращает middleware для доступа
func (p *provider) Access() *access.Middleware {
	return p.accessMiddleware
}

// Metrics возвращает middleware для метрик
func (p *provider) Metrics() *metrics.Middleware {
	return p.metricsMiddleware
}

// TimeoutMiddleware ограничивает выполнение обработчика по времени
func (p *provider) TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		ch := make(chan struct{})

		c.Request = c.Request.WithContext(ctx)

		go func() {
			c.Next()
			close(ch)
		}()

		select {
		case <-ch:
			return
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "request timeout"})
		}
	}
}
