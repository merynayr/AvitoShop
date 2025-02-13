package access

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/utils/jwt"
)

const (
	authHeader = "Authorization"
	authPrefix = "Bearer "
)

// Check проверяет, имеет ли пользователь доступ к эндпоинту
func (s *srv) Check(ctx *gin.Context, endpointAddress string) (string, error) {

	if _, ok := s.userAccesses[endpointAddress]; !ok {
		return "", nil
	}

	authHeader := ctx.GetHeader(authHeader)
	if authHeader == "" {
		return "", errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader, authPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader, authPrefix)

	claims, err := jwt.VerifyToken(accessToken, s.authConfig.AccessTokenSecretKey())
	if err != nil {
		return "", errors.New("access token is invalid")
	}

	// Проверяем, есть ли доступ у пользователя к запрашиваемому пути
	if _, ok := s.userAccesses[endpointAddress]; !ok {
		return "", errors.New("access denied")
	}

	return claims.Username, nil
}
