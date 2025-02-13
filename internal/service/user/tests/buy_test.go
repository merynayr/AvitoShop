package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/merynayr/AvitoShop/internal/client/db"
	txMocks "github.com/merynayr/AvitoShop/internal/client/db/mocks"
	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/repository"
	repositoryMocks "github.com/merynayr/AvitoShop/internal/repository/mocks"
	"github.com/merynayr/AvitoShop/internal/service/user"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

func TestBuy(t *testing.T) {
	t.Parallel()

	type shopRepositoryMockFunc func(mc *minimock.Controller) repository.ShopRepository
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		user *model.User
		item string
	}

	var (
		ctx = context.Background()

		userID    = gofakeit.Int64()
		userCoins = int64(100)
		item      = "item1"
		itemPrice = int64(50)
		Quantity  = int64(1)
		repoErr   = fmt.Errorf("repository error")

		userModel = &model.User{
			ID:    userID,
			Coins: userCoins,
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
				ctx:  ctx,
				user: userModel,
				item: item,
			},
			err: nil,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(itemPrice, nil)
				mock.CheckInventoryMock.Expect(ctx, userID, item).Return(true, Quantity, nil)
				mock.UpdateInventoryMock.Expect(ctx, item, userID, Quantity+1).Return(nil)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, &model.UserUpdate{
					ID:    userID,
					Coins: userCoins - itemPrice,
				}).Return(nil)
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
			name: "not enough coins",
			args: args{
				ctx:  ctx,
				user: userModel,
				item: item,
			},
			err: sys.NewCommonError("not enough coins", codes.BadRequest),
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(500, nil)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)

				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "item not found",
			args: args{
				ctx:  ctx,
				item: item,
			},
			err: sys.NewCommonError("item not found", codes.NotFound),
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(0, repoErr)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
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
				ctx:  ctx,
				user: userModel,
				item: item,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(itemPrice, nil)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
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
			name: "error in UpdateUser",
			args: args{
				ctx:  ctx,
				user: userModel,
				item: item,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(itemPrice, nil)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, &model.UserUpdate{
					ID:    userID,
					Coins: userModel.Coins - itemPrice,
				}).Return(repoErr)
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
			name: "error in CheckInventory",
			args: args{
				ctx:  ctx,
				user: userModel,
				item: item,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(itemPrice, nil)
				mock.CheckInventoryMock.Expect(ctx, userID, item).Return(false, 0, repoErr)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, &model.UserUpdate{
					ID:    userID,
					Coins: userModel.Coins - itemPrice,
				}).Return(nil)
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
			name: "error in InsertNewInventory",
			args: args{
				ctx:  ctx,
				user: userModel,
				item: item,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(itemPrice, nil)
				mock.CheckInventoryMock.Expect(ctx, userID, item).Return(false, 0, nil)
				mock.InsertNewInventoryMock.Expect(ctx, userID, item).Return(repoErr)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, &model.UserUpdate{
					ID:    userID,
					Coins: userModel.Coins - itemPrice,
				}).Return(nil)
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
			name: "error in UpdateInventory",
			args: args{
				ctx:  ctx,
				user: userModel,
				item: item,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, item).Return(itemPrice, nil)
				mock.CheckInventoryMock.Expect(ctx, userID, item).Return(true, Quantity, nil)
				mock.UpdateInventoryMock.Expect(ctx, item, userID, Quantity+1).Return(repoErr)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, &model.UserUpdate{
					ID:    userID,
					Coins: userModel.Coins - itemPrice,
				}).Return(nil)
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
		logger.Init("debug")
		t.Run(tt.name, func(t *testing.T) {
			shopRepoMock := tt.shopRepositoryMock(minimock.NewController(t))
			userRepoMock := tt.userRepositoryMock(minimock.NewController(t))
			txManagerMock := tt.txManagerMock(minimock.NewController(t))

			service := user.NewService(shopRepoMock, userRepoMock, txManagerMock)
			err := service.Buy(tt.args.ctx, tt.args.user, tt.args.item)

			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
