package user

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

func (s *userService) SendCoins(ctx context.Context, fromUser *model.User, sendCoins *model.SendCoinRequest) error {
	toUser, err := s.userRepository.GetUserByName(ctx, sendCoins.ToUser)
	if err != nil {
		return sys.NewCommonError("recipient not found", codes.BadRequest)
	}

	if fromUser.Coins < sendCoins.Amount {
		return sys.NewCommonError("not enough coins", codes.BadRequest)
	}

	if toUser.ID == fromUser.ID {
		return sys.NewCommonError("you can't transfer money to yourself", codes.BadRequest)
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    fromUser.ID,
			Coins: fromUser.Coins - sendCoins.Amount,
		})
		if errTx != nil {
			return errTx
		}
		errTx = s.userRepository.UpdateUser(ctx, &model.UserUpdate{
			ID:    toUser.ID,
			Coins: toUser.Coins + sendCoins.Amount,
		})
		if errTx != nil {
			return errTx
		}

		errTx = s.shopRepository.CreateTransaction(ctx, fromUser.ID, toUser.ID, sendCoins.Amount)
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		logger.Error(err.Error())
	}
	return err
}
