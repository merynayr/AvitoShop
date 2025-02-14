package access

import (
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"
)

type srv struct {
	userService  service.UserService
	userAccesses map[string]struct{}
	authConfig   config.AuthConfig
}

// NewService возвращает новый объект сервисного слоя access
func NewService(userService service.UserService, userAccesses map[string]struct{}, authConfig config.AuthConfig) service.AccessService {
	return &srv{
		userService:  userService,
		userAccesses: userAccesses,
		authConfig:   authConfig,
	}
}
