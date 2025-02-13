package middleware

import "github.com/gin-gonic/gin"

// UserMiddlware интерфейс
type UserMiddlware interface {
	ExtractUserID() gin.HandlerFunc
}
