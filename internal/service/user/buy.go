package user

import (
	"context"
	"fmt"

	"github.com/merynayr/AvitoShop/internal/model"
)

func (s *userService) Buy(ctx context.Context, user *model.User, item string) error {
	price, err := s.shopRepository.GetMerchPrice(ctx, item)
	if err != nil {
		return err
	}

	if user.Coins < price {
		return fmt.Errorf("errors: not enough coins")
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

		exist, newQuantity, errTx := s.shopRepository.CheckInventory(ctx, user, item)
		if errTx != nil {
			return errTx
		}

		if !exist {
			errTx = s.shopRepository.InsertNewInventory(ctx, user, item)
			if errTx != nil {
				return errTx
			}
		} else {
			errTx = s.shopRepository.UpdateInventory(ctx, user.ID, newQuantity)
			if errTx != nil {
				return errTx
			}
		}

		return nil
	})

	return err
}
