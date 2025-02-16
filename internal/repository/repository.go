package repository

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/model"
)

// ShopRepository - интерфейс репо слоя shop
type ShopRepository interface {
	GetMerchPrice(ctx context.Context, item string) (int64, error)
	CheckInventory(ctx context.Context, userID int64, item string) (bool, int64, error)
	InsertNewInventory(ctx context.Context, userID int64, item string) error
	UpdateInventory(ctx context.Context, item string, id, newQuantity int64) error
	CreateTransaction(ctx context.Context, fromUserID, toUserID, amount int64) error
	GetUserTransactions(ctx context.Context, userID int64) ([]model.Received, []model.Sent, error)
	GetUserInventory(ctx context.Context, userID int64) ([]model.InventoryItem, error)
}

// UserRepository - интерфейс репо слоя user
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) error
	GetUserByName(ctx context.Context, name string) (*model.User, error)
}
