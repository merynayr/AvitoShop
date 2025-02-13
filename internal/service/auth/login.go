package auth

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
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

	user, err := s.userRepository.GetUserByName(ctx, name)
	if err != nil {
		return "", err
	}
	userInfo := &model.AuthRequest{
		Username: user.Username,
		Password: user.Password,
	}

	err = hash.CompareHashAndPass(password, userInfo.Password)
	if err != nil {
		return "", sys.NewCommonError("invalid password", codes.BadRequest)
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.RefreshTokenSecretKey(), s.authCfg.RefreshTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}
