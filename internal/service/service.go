package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
)

// ShopService интерфейс сервисного слоя shop
type ShopService interface {
	GetMerchPrice(ctx context.Context, item string) (int64, error)
}

// UserService интерфейс сервисного слоя user
type UserService interface {
	Buy(ctx context.Context, user *model.User, item string) error
	GetUserByName(ctx context.Context, name string) (*model.User, error)
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
