package metric

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// name
const (
	namespace = "shop_space"
	appName   = "shop_app"
	subsystem = "http"
)

// Metrics структура объекта сборщика метрик
type Metrics struct {
	requestCounter        *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
	responseCounter       *prometheus.CounterVec
	successfulResponses   *prometheus.CounterVec
	failedResponses       *prometheus.CounterVec
}

var metrics *Metrics

// Init инициализирует сборщик метрик
func Init(_ context.Context) {
	metrics = &Metrics{
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_requests_total",
				Help:      "Общее количество HTTP-запросов",
			},
			[]string{"method", "path"},
		),
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_responses_total",
				Help:      "Общее количество HTTP-ответов",
			},
			[]string{"status", "method", "path"},
		),
		successfulResponses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_successful_responses_total",
				Help:      "Общее количество успешных HTTP-ответов (2xx)",
			},
			[]string{"method", "path"},
		),
		failedResponses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_failed_responses_total",
				Help:      "Общее количество ошибочных HTTP-ответов (4xx и 5xx)",
			},
			[]string{"method", "path"},
		),
		histogramResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_response_time_seconds",
				Help:      "Гистограмма времени ответа HTTP",
				Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
			},
			[]string{"status", "method"},
		),
	}
}

// IncRequestCounter увеличивает requestCounter
func IncRequestCounter(method, path string) {
	metrics.requestCounter.WithLabelValues(method, path).Inc()
}

// HistogramResponseTimeObserve сохраняет время ответа
func HistogramResponseTimeObserve(status, method string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status, method).Observe(time)
}

// IncResponseCounter увеличивает responseCounter и считает успешные/ошибочные ответы
func IncResponseCounter(status, method, path string) {
	metrics.responseCounter.WithLabelValues(status, method, path).Inc()

	switch status {
	case "success":
		metrics.successfulResponses.WithLabelValues(method, path).Inc()
	case "error":
		metrics.failedResponses.WithLabelValues(method, path).Inc()
	}
}
