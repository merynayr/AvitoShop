package repository

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/model"
)

// ShopRepository - интерфейс репо слоя shop
type ShopRepository interface {
}

// UserRepository - интерфейс репо слоя user
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) error
	DeleteUser(ctx context.Context, userID int64) error
	GetUserByName(ctx context.Context, name string) (*model.UserInfo, error)
	IsNameExist(ctx context.Context, name string) (bool, error)
}
