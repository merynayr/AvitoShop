package auth

import (
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/repository"
	"github.com/merynayr/AvitoShop/internal/service"
)

type srv struct {
	userRepository repository.UserRepository
	authCfg        config.AuthConfig
}

// NewService возвращает новый объект сервисного слоя AvitoShop
func NewService(userRepo repository.UserRepository, authCfg config.AuthConfig) service.AuthService {
	return &srv{
		userRepository: userRepo,
		authCfg:        authCfg,
	}
}
