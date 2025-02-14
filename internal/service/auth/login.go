package auth

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/utils/hash"
	"github.com/merynayr/AvitoShop/internal/utils/jwt"
)

// Login валидирует данные пользователя, и если все ок, возвращает token-ы
func (s *srv) Login(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	var userInfo *model.AuthRequest
	user, err := s.userRepository.GetUserByName(ctx, username)
	if err != nil {
		// Если пользователя не существует, то создаём его
		if sys.GetCommonError(err) == sys.UserNotFoundError {
			_, err := s.userRepository.CreateUser(ctx, &model.User{
				Username: username,
				Password: password,
				Coins:    1000,
			})
			if err != nil {
				return nil, err
			}
			userInfo = &model.AuthRequest{
				Username: username,
				Password: password,
			}
		} else {
			return nil, err
		}
	} else {
		userInfo = &model.AuthRequest{
			Username: user.Username,
			Password: user.Password,
		}

		err = hash.CompareHashAndPass(password, userInfo.Password)
		if err != nil {
			return nil, sys.InvalidPasswordError
		}
	}

	refreshToken, err := jwt.GenerateToken(userInfo, s.authCfg.RefreshTokenSecretKey(), s.authCfg.RefreshTokenExp())
	if err != nil {
		return nil, err
	}
	accessToken, err := jwt.GenerateToken(userInfo, s.authCfg.AccessTokenSecretKey(), s.authCfg.AccessTokenExp())
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
