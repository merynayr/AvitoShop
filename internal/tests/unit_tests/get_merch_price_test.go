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
	"github.com/merynayr/AvitoShop/internal/repository"
	repositoryMocks "github.com/merynayr/AvitoShop/internal/repository/mocks"
	"github.com/merynayr/AvitoShop/internal/service/shop"
)

func TestGetMerchPrice(t *testing.T) {
	t.Parallel()

	type shopRepositoryMockFunc func(mc *minimock.Controller) repository.ShopRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx = context.Background()

		item = gofakeit.Username()

		req     = item
		price   = int64(100)
		repoErr = fmt.Errorf("repo error")
	)
	tests := []struct {
		item               string
		args               args
		err                error
		want               int64
		shopRepositoryMock shopRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			item: "success from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, req).Return(price, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			item: "error from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			shopRepositoryMock: func(mc *minimock.Controller) repository.ShopRepository {
				mock := repositoryMocks.NewShopRepositoryMock(mc)
				mock.GetMerchPriceMock.Expect(ctx, req).Return(0, repoErr)
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
		t.Run(tt.item, func(t *testing.T) {
			shopRepoMock := tt.shopRepositoryMock(minimock.NewController(t))
			txManagerMock := tt.txManagerMock(minimock.NewController(t))

			service := shop.NewService(shopRepoMock, nil, txManagerMock)
			_, err := service.GetMerchPrice(tt.args.ctx, tt.args.req)

			if tt.err != nil {
				require.Equal(t, tt.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
