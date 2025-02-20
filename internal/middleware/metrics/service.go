package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/metric"
)

// Middleware метрик структура
type Middleware struct {
}

// NewMetricsMiddleware возвращает новый объект middleware слоя метрик
func NewMetricsMiddleware() *Middleware {
	return &Middleware{}
}

// Metrics сохраняет метрики для каждого запроса
func (m *Middleware) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		metric.IncRequestCounter(c.Request.Method, path)
		timeStart := time.Now()

		c.Next()

		diffTime := time.Since(timeStart)
		status := "success"
		if c.Writer.Status() >= 400 {
			status = "error"
		}

		metric.IncResponseCounter(status, c.Request.Method, path)
		metric.HistogramResponseTimeObserve(status, c.Request.Method, diffTime.Seconds())
	}
}
