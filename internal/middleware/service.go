package middleware

import (
	"github.com/gin-gonic/gin"
)

// Middleware интерфейс
type Middleware interface {
	Check() gin.HandlerFunc
}
