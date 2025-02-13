package user

import (
	"context"
	"errors"

	"github.com/merynayr/AvitoShop/internal/model"
)

func (s *userService) SendCoins(ctx context.Context, userID int64, SendCoins model.SendCoinRequest) error {
	toUser, err := s.userRepository.GetUserByName(ctx, SendCoins.ToUser)
	if err != nil {
		return errors.New("recipient not found")
	}

	fromUser, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if fromUser.Coins < SendCoins.Amount {
		return errors.New("not enough coins")
	}

	if toUser.ID == fromUser.ID {
		return errors.New("you can't transfer money to yourself")
	}
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    fromUser.ID,
			Coins: fromUser.Coins - SendCoins.Amount,
		})
		if errTx != nil {
			return errTx
		}
		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    toUser.ID,
			Coins: toUser.Coins + SendCoins.Amount,
		})
		if errTx != nil {
			return errTx
		}

		errTx = s.shopRepository.CreateTransaction(ctx, fromUser.ID, toUser.ID, SendCoins.Amount)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	return err
}
