package user

import (
	"context"
	"errors"

	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
)

func (s *userService) GetUserInfo(ctx context.Context, userID int64) (model.UserInfoResponse, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return model.UserInfoResponse{}, errors.New("user not found")
	}

	items, err := s.shopRepository.GetUserInventory(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return model.UserInfoResponse{}, errors.New("cannot get inventory")
	}

	received, sent, err := s.shopRepository.GetUserTransactions(ctx, userID)
	if err != nil {
		logger.Error(err.Error())
		return model.UserInfoResponse{}, errors.New("cannot get transactions")
	}

	return model.UserInfoResponse{
		Coins:     user.Coins,
		Inventory: items,
		CoinHistory: model.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
