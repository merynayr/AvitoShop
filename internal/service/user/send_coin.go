package user

import (
	"context"
	"errors"

	"github.com/merynayr/AvitoShop/internal/model"
)

func (s *userService) SendCoins(ctx context.Context, userID int64, toUsername string, amount int64) error {
	toUser, err := s.userRepository.GetUserByName(ctx, toUsername)
	if err != nil {
		return errors.New("recipient not found")
	}

	fromUser, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if fromUser.Coins < amount {
		return errors.New("not enough coins")
	}

	if toUser.ID == fromUser.ID {
		return errors.New("you can't transfer money to yourself")
	}
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    fromUser.ID,
			Coins: fromUser.Coins - amount,
		})
		if errTx != nil {
			return errTx
		}
		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    toUser.ID,
			Coins: toUser.Coins + amount,
		})
		if errTx != nil {
			return errTx
		}

		errTx = s.shopRepository.CreateTransaction(ctx, fromUser.ID, toUser.ID, amount)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	return err
}
