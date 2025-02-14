package auth

import (
	"context"
	"fmt"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/utils/jwt"
)

// GetRefreshToken принимает старый refresh token и на его основе создает и возвращает новый
func (s *srv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	claims, err := jwt.VerifyToken(oldRefreshToken, s.authCfg.RefreshTokenSecretKey())
	if err != nil {
		return "", fmt.Errorf(sys.ErrInvalidRefreshToken)
	}

	user, err := s.userRepository.GetUserByName(ctx, claims.Username)
	if err != nil {
		return "", err
	}
	userInfo := &model.AuthRequest{
		Username: user.Username,
		Password: user.Password,
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.RefreshTokenSecretKey(), s.authCfg.RefreshTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}
