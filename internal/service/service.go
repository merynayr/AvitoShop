package service

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ShopService интерфейс сервисного слоя shop
type ShopService interface {
}

// UserService интерфейс сервисного слоя user
type UserService interface {
}

// AuthService интерфейс сервисного слоя auth
type AuthService interface {
	Login(ctx context.Context, name string, password string) (string, error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// AccessService интерфейс сервисного слоя access
type AccessService interface {
	Check(ctx *gin.Context, endpointAddress string) (string, error)
}
