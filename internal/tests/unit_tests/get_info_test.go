package unittests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/merynayr/AvitoShop/internal/client/db"
	txMocks "github.com/merynayr/AvitoShop/internal/client/db/mocks"
	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/repository"
	repositoryMocks "github.com/merynayr/AvitoShop/internal/repository/mocks"
	"github.com/merynayr/AvitoShop/internal/service/shop"
)

func TestGetInfo(t *testing.T) {
	t.Parallel()

	type shopRepositoryMockFunc func(mc *minimock.Controller) repository.ShopRepository
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		user *model.User
	}

	var (
		ctx = context.Background()

		userID    = gofakeit.Int64()
		userCoins = int64(100)

		repoErr = fmt.Errorf("repository error")

		userModel = &model.User{
			ID:    userID,
			Coins: userCoins,
		}

		items = []model.InventoryItem{{
			ItemName: "Item",
			Quantity: 3,
		},
		}

		received = []model.Received{{
			ID:           1,
			FromUsername: "Oleg",
			Amount:       500,
			CreatedAt:    time.Now(),
		}}

		sent = []model.Sent{
			{
				ID:         17,
				ToUsername: "Admin",
				Amount:     200,
				CreatedAt:  time.Now(),
			},
		}
	)

	tests := []struct {
		name               string
		args               args
		err                error
		want               model.InfoResponse
		shopRepositoryMock shopRepositoryMockFunc
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				user: userModel,
			},
			err: nil,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetUserInventoryMock.Expect(ctx, userID).Return(items, nil)
				mock.GetUserTransactionsMock.Expect(ctx, userID).Return(received, sent, nil)
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
			name: "error in GetUserInventory",
			args: args{
				ctx:  ctx,
				user: userModel,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetUserInventoryMock.Expect(ctx, userID).Return(nil, repoErr)
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
			name: "error in GetUserTransactions",
			args: args{
				ctx:  ctx,
				user: userModel,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetUserInventoryMock.Expect(ctx, userID).Return(items, nil)
				mock.GetUserTransactionsMock.Expect(ctx, userID).Return(nil, nil, repoErr)
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
	}

	for _, tt := range tests {
		tt := tt
		logger.Init("debug")
		t.Run(tt.name, func(t *testing.T) {
			shopRepoMock := tt.shopRepositoryMock(minimock.NewController(t))
			userRepoMock := tt.userRepositoryMock(minimock.NewController(t))
			txManagerMock := tt.txManagerMock(minimock.NewController(t))

			service := shop.NewService(shopRepoMock, userRepoMock, txManagerMock)
			_, err := service.GetUserInfo(tt.args.ctx, tt.args.user)

			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
