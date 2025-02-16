package unittests

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
	"github.com/merynayr/AvitoShop/internal/service/shop"
)

func TestGetUserByName(t *testing.T) {
	t.Parallel()

	type shopRepositoryMockFunc func(mc *minimock.Controller) repository.ShopRepository
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx = context.Background()

		name = gofakeit.Username()

		req = name

		userModel = &model.User{}

		repoErr = fmt.Errorf("repo error")
	)
	tests := []struct {
		name               string
		args               args
		err                error
		want               model.User
		shopRepositoryMock shopRepositoryMockFunc
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: *userModel,
			err:  nil,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, req).Return(userModel, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "error from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				return mock
			},
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserByNameMock.Expect(ctx, req).Return(nil, repoErr)
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
			_, err := service.GetUserByName(tt.args.ctx, tt.args.req)

			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
