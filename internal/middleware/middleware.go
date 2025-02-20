package middleware

import (
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/middleware/access"
	"github.com/merynayr/AvitoShop/internal/middleware/metrics"
	"github.com/merynayr/AvitoShop/internal/service"
)

// Middleware интерфейс для всех middleware
type Middleware interface {
	Access() *access.Middleware
	Metrics() *metrics.Middleware
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
