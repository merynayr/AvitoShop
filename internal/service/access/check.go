package access

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/utils/jwt"
)

const (
	authHeader = "Authorization"
	authPrefix = "Bearer "
)

// Check проверяет, имеет ли пользователь доступ к эндпоинту
func (s *srv) Check(ctx *gin.Context, endpointAddress string) (*model.User, error) {

	if _, ok := s.userAccesses[endpointAddress]; !ok {
		return nil, nil
	}

	authHeader := ctx.GetHeader(authHeader)
	if authHeader == "" {
		return nil, errors.New(sys.ErrAuthHeaderNotProvided)
	}

	if !strings.HasPrefix(authHeader, authPrefix) {
		return nil, errors.New(sys.ErrInvalidAuthHeaderFormat)
	}

	accessToken := strings.TrimPrefix(authHeader, authPrefix)

	claims, err := jwt.VerifyToken(accessToken, s.authConfig.AccessTokenSecretKey())
	if err != nil {
		return nil, errors.New(sys.ErrInvalidAccessToken)
	}

	user, err := s.shopService.GetUserByName(ctx, claims.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
