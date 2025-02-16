package shop

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
)

func (s *shopService) Buy(ctx context.Context, user *model.User, item string) error {
	price, err := s.shopRepository.GetMerchPrice(ctx, item)
	if err != nil {
		return sys.ItemNotFoundError
	}

	if user.Coins < price {
		return sys.NotEnoughCoinsError
	}
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    user.ID,
			Coins: user.Coins - price,
		})
		if errTx != nil {
			return errTx
		}

		exist, Quantity, errTx := s.shopRepository.CheckInventory(ctx, user.ID, item)
		if errTx != nil {
			return errTx
		}

		if !exist {
			errTx = s.shopRepository.InsertNewInventory(ctx, user.ID, item)
			if errTx != nil {
				return errTx
			}
		} else {
			errTx = s.shopRepository.UpdateInventory(ctx, item, user.ID, Quantity+1)
			if errTx != nil {
				return errTx
			}
		}

		return nil
	})
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}
