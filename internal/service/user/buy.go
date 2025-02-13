package user

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

func (s *userService) Buy(ctx context.Context, user *model.User, item string) error {
	price, err := s.shopRepository.GetMerchPrice(ctx, item)
	if err != nil {
		return sys.NewCommonError("item does not exist", codes.BadRequest)
	}

	if user.Coins < price {
		return sys.NewCommonError("not enough coins", codes.BadRequest)
	}
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    user.ID,
			Coins: user.Coins - price,
		})
		if errTx != nil {
			logger.Error(errTx.Error())
			return errTx
		}

		exist, Quantity, errTx := s.shopRepository.CheckInventory(ctx, user.ID, item)
		if errTx != nil {
			logger.Error(errTx.Error())
			return errTx
		}

		if !exist {
			errTx = s.shopRepository.InsertNewInventory(ctx, user.ID, item)
			if errTx != nil {
				logger.Error(errTx.Error())
				return errTx
			}
		} else {
			errTx = s.shopRepository.UpdateInventory(ctx, item, user.ID, Quantity+1)
			if errTx != nil {
				logger.Error(errTx.Error())
				return errTx
			}
		}

		return nil
	})

	return err
}
