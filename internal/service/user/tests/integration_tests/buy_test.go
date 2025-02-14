package integrationtests

import (
	"context"
	"fmt"
	"testing"

	"github.com/merynayr/AvitoShop/internal/model"
)

func (s *MyNewIntegrationSuite) TestBuy() {
	ctx := context.Background()

	userSrv := s.SetupTest()

	user := &model.User{
		ID:       1,
		Username: "TestUser",
		Password: "1",
		Coins:    100,
	}

	tests := []struct {
		name            string
		product         string
		expectedError   error
		expectedBalance int64
	}{
		{
			name:            "Purchase cup",
			product:         "cup",
			expectedError:   nil,
			expectedBalance: user.Coins - merchPrices["cup"],
		},
		{
			name:            "Purchase book",
			product:         "book",
			expectedError:   nil,
			expectedBalance: user.Coins - merchPrices["book"],
		},
		{
			name:            "not enough coins",
			product:         "pink-hoody",
			expectedError:   fmt.Errorf("not enough coins"),
			expectedBalance: user.Coins,
		},
		{
			name:            "item not found",
			product:         "hat",
			expectedError:   fmt.Errorf("item not found"),
			expectedBalance: user.Coins,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(_ *testing.T) {

			_, err := s.pool.Exec(ctx, "INSERT INTO users (id, username, password, coins) VALUES ($1, $2, $3, $4)", user.ID, user.Username, user.Password, user.Coins)
			s.Require().NoError(err)

			// Выполняем тестируемую функцию
			err = userSrv.Buy(ctx, user, tt.product)

			if tt.expectedError != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tt.expectedError.Error())
			} else {
				s.Require().NoError(err)
			}

			var newBalance int64
			err = s.pool.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1", user.ID).Scan(&newBalance)
			s.Require().NoError(err)
			s.Require().Equal(tt.expectedBalance, newBalance, "Баланс пользователя ID %d не соответствует ожидаемому", user.ID)

			s.TearDownTest()
		})
	}
}
