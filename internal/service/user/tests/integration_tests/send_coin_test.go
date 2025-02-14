package integrationtests

import (
	"context"
	"testing"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
)

func (s *MyNewIntegrationSuite) TestSendCoins() {
	ctx := context.Background()
	userSrv := s.SetupTest()

	tests := []struct {
		name            string
		fromUser        *model.User
		toUser          *model.User
		sendCoins       *model.SendCoinRequest
		setupMock       func()
		expectedError   error
		expectedBalance map[int64]int64
	}{
		{
			name:     "Successful transfer",
			fromUser: &model.User{ID: 1, Username: "Alice", Password: "1", Coins: 100},
			toUser:   &model.User{ID: 2, Username: "Bob", Password: "1", Coins: 50},
			sendCoins: &model.SendCoinRequest{
				ToUser: "Bob",
				Amount: 30,
			},
			expectedError: nil,
			expectedBalance: map[int64]int64{
				1: 70, // Alice (100 - 30)
				2: 80, // Bob (50 + 30)
			},
		},
		{
			name:     "Not enough coins",
			fromUser: &model.User{ID: 1, Username: "Alice", Password: "1", Coins: 20},
			toUser:   &model.User{ID: 2, Username: "Bob", Password: "1", Coins: 50},
			sendCoins: &model.SendCoinRequest{
				ToUser: "Bob",
				Amount: 50,
			},
			expectedError: sys.NotEnoughCoinsError,
			expectedBalance: map[int64]int64{
				1: 20,
				2: 50,
			},
		},
		{
			name:     "Transfer to self",
			fromUser: &model.User{ID: 1, Username: "Alice", Password: "1", Coins: 100},
			sendCoins: &model.SendCoinRequest{
				ToUser: "Alice",
				Amount: 10,
			},
			expectedError: sys.SelfTransferNotAllowedError,
			expectedBalance: map[int64]int64{
				1: 100,
			},
		},
		{
			name:     "Recipient not found",
			fromUser: &model.User{ID: 1, Username: "Alice", Password: "1", Coins: 100},
			toUser:   nil,
			sendCoins: &model.SendCoinRequest{
				ToUser: "Charlie",
				Amount: 10,
			},
			expectedError: sys.RecipientNotFoundError,
			expectedBalance: map[int64]int64{
				1: 100,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(_ *testing.T) {

			if tt.fromUser != nil {
				_, err := s.pool.Exec(ctx, "INSERT INTO users (id, username, password, coins) VALUES ($1, $2, $3, $4)",
					tt.fromUser.ID, tt.fromUser.Username, tt.fromUser.Password, tt.fromUser.Coins)
				s.Require().NoError(err)
			}
			if tt.toUser != nil {
				_, err := s.pool.Exec(ctx, "INSERT INTO users (id, username, password, coins) VALUES ($1, $2, $3, $4)",
					tt.toUser.ID, tt.toUser.Username, tt.toUser.Password, tt.toUser.Coins)
				s.Require().NoError(err)
			}

			// Выполняем тестируемую функцию
			err := userSrv.SendCoins(ctx, tt.fromUser, tt.sendCoins)

			if tt.expectedError != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tt.expectedError.Error())
			} else {
				s.Require().NoError(err)
			}

			for userID, expectedBalance := range tt.expectedBalance {
				var actualBalance int64
				err := s.pool.QueryRow(ctx, "SELECT coins FROM users WHERE id = $1", userID).Scan(&actualBalance)
				s.Require().NoError(err)
				s.Require().Equal(expectedBalance, actualBalance, "Баланс пользователя ID %d не соответствует ожидаемому", userID)
			}

			s.TearDownTest()
		})
	}
}
