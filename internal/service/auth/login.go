package auth

import (
	"context"
	"fmt"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/utils/hash"
	"github.com/merynayr/AvitoShop/internal/utils/jwt"
)

// Login валидирует данные пользователя, и если все ок, возвращает refresh token
func (s *srv) Login(ctx context.Context, name string, password string) (string, error) {
	exist, err := s.userRepository.IsNameExist(ctx, name)
	if err != nil {
		return "", err
	}
	if !exist {
		_, err := s.userRepository.CreateUser(ctx, &model.User{
			Username: name,
			Password: password,
			Coins:    1000,
		})
		if err != nil {
			return "", err
		}
	}

	userInfo, err := s.userRepository.GetUserByName(ctx, name)
	if err != nil {
		return "", err
	}

	err = hash.CompareHashAndPass(password, userInfo.Password)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.RefreshTokenSecretKey(), s.authCfg.RefreshTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}
