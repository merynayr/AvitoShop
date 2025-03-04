package shop

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/model"
)

func (s *shopService) GetUserInfo(ctx context.Context, user *model.User) (*model.InfoResponse, error) {
	items, err := s.shopRepository.GetUserInventory(ctx, user.ID)
	if err != nil {
		return &model.InfoResponse{}, err
	}

	received, sent, err := s.shopRepository.GetUserTransactions(ctx, user.ID)
	if err != nil {
		return &model.InfoResponse{}, err
	}

	return &model.InfoResponse{
		Coins:     user.Coins,
		Inventory: items,
		CoinHistory: model.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
