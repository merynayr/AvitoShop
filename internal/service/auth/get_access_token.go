package auth

import (
	"context"
	"fmt"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/utils/jwt"
)

// GetAccessToken принимает refresh token и на его основе создает и возвращает access token
func (s *srv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := jwt.VerifyToken(refreshToken, s.authCfg.RefreshTokenSecretKey())
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	user, err := s.userRepository.GetUserByName(ctx, claims.Username)
	if err != nil {
		return "", err
	}
	userInfo := &model.UserInfo{
		Username: user.Username,
		Password: user.Password,
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.AccessTokenSecretKey(), s.authCfg.AccessTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}
