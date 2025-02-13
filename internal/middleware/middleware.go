package middleware

import "github.com/gin-gonic/gin"

// UserMiddleware интерфейс
type UserMiddleware interface {
	AddAccessTokenFromCookie() gin.HandlerFunc
	Check() gin.HandlerFunc
}
