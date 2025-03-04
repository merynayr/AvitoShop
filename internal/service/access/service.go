package access

import (
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"
)

type srv struct {
	shopService  service.ShopService
	userAccesses map[string]struct{}
	authConfig   config.AuthConfig
}

// NewService возвращает новый объект сервисного слоя access
func NewService(shopService service.ShopService, userAccesses map[string]struct{}, authConfig config.AuthConfig) service.AccessService {
	return &srv{
		shopService:  shopService,
		userAccesses: userAccesses,
		authConfig:   authConfig,
	}
}
