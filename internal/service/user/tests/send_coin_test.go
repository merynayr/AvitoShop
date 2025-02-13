package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/repository"
	"github.com/merynayr/AvitoShop/internal/service/user"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"

	txMocks "github.com/merynayr/AvitoShop/internal/client/db/mocks"
	repositoryMocks "github.com/merynayr/AvitoShop/internal/repository/mocks"
)

func TestSendCoin(t *testing.T) {
	t.Parallel()

	type shopRepositoryMockFunc func(mc *minimock.Controller) repository.ShopRepository
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx       context.Context
		fromUser  *model.User
		SendCoins *model.SendCoinRequest
	}

	var (
		ctx = context.Background()

		FromUserID = int64(1)
		ToUserID   = int64(2)

		repoErr = fmt.Errorf("repository error")

		fromUser = &model.User{
			ID:    FromUserID,
			Coins: int64(100),
		}

		toUser = &model.User{
			ID:    ToUserID,
			Coins: int64(100),
		}

		sendCoins = &model.SendCoinRequest{
			ToUser: "Admin",
			Amount: 20,
		}
	)

	tests := []struct {
		name               string
		args               args
		err                error
		shopRepositoryMock shopRepositoryMockFunc
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:       ctx,
				fromUser:  fromUser,
				SendCoins: sendCoins,
			},
			err: nil,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.CreateTransactionMock.Expect(ctx, fromUser.ID, toUser.ID, sendCoins.Amount).Return(nil)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)

				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(toUser, nil)

				mock.UpdateUserMock.
					When(ctx, &model.UserUpdate{
						ID:    FromUserID,
						Coins: fromUser.Coins - sendCoins.Amount,
					}).
					Then(nil)

				mock.UpdateUserMock.
					When(ctx, &model.UserUpdate{
						ID:    ToUserID,
						Coins: toUser.Coins + sendCoins.Amount,
					}).
					Then(nil)

				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "recipient not found",
			args: args{
				ctx:       ctx,
				fromUser:  fromUser,
				SendCoins: sendCoins,
			},
			err: sys.NewCommonError("recipient not found", codes.BadRequest),
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(nil, repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "not enough coins",
			args: args{
				ctx:      ctx,
				fromUser: fromUser,
				SendCoins: &model.SendCoinRequest{
					ToUser: "Admin",
					Amount: 120,
				},
			},
			err: sys.NewCommonError("not enough coins", codes.BadRequest),
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(toUser, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "you can't transfer money to yourself",
			args: args{
				ctx: ctx,
				fromUser: &model.User{
					ID:    2,
					Coins: int64(100),
				},
				SendCoins: sendCoins,
			},
			err: sys.NewCommonError("you can't transfer money to yourself", codes.BadRequest),
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(toUser, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "transaction error",
			args: args{
				ctx:       ctx,
				fromUser:  fromUser,
				SendCoins: sendCoins,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)

				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(toUser, nil)

				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(_ context.Context, _ db.Handler) (err error) {
					return repoErr
				})
				return mock
			},
		},
		{
			name: "error in UpdateUser1",
			args: args{
				ctx:       ctx,
				fromUser:  fromUser,
				SendCoins: sendCoins,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(toUser, nil)

				mock.UpdateUserMock.
					When(ctx, &model.UserUpdate{
						ID:    FromUserID,
						Coins: fromUser.Coins - sendCoins.Amount,
					}).
					Then(repoErr)

				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "error in UpdateUser2",
			args: args{
				ctx:       ctx,
				fromUser:  fromUser,
				SendCoins: sendCoins,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(toUser, nil)

				mock.UpdateUserMock.
					When(ctx, &model.UserUpdate{
						ID:    FromUserID,
						Coins: fromUser.Coins - sendCoins.Amount,
					}).
					Then(nil)

				mock.UpdateUserMock.
					When(ctx, &model.UserUpdate{
						ID:    toUser.ID,
						Coins: toUser.Coins + sendCoins.Amount,
					}).
					Then(repoErr)

				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "error in CreateTransaction",
			args: args{
				ctx:       ctx,
				fromUser:  fromUser,
				SendCoins: sendCoins,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.CreateTransactionMock.Expect(ctx, fromUser.ID, toUser.ID, sendCoins.Amount).Return(repoErr)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, sendCoins.ToUser).Return(toUser, nil)

				mock.UpdateUserMock.
					When(ctx, &model.UserUpdate{
						ID:    FromUserID,
						Coins: fromUser.Coins - sendCoins.Amount,
					}).
					Then(nil)

				mock.UpdateUserMock.
					When(ctx, &model.UserUpdate{
						ID:    toUser.ID,
						Coins: toUser.Coins + sendCoins.Amount,
					}).
					Then(nil)

				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		logger.Init("info")
		t.Run(tt.name, func(t *testing.T) {
			shopRepoMock := tt.shopRepositoryMock(minimock.NewController(t))
			userRepoMock := tt.userRepositoryMock(minimock.NewController(t))
			txManagerMock := tt.txManagerMock(minimock.NewController(t))

			service := user.NewService(shopRepoMock, userRepoMock, txManagerMock)
			err := service.SendCoins(tt.args.ctx, tt.args.fromUser, tt.args.SendCoins)

			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
