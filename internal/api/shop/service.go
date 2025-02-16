package shop

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/service"
)

// API user структура
// объект сервисного слоя (его интерфейса)
type API struct {
	shopService service.ShopService
}

// NewAPI возвращает новый объект имплементации API-слоя
func NewAPI(shopService service.ShopService) *API {
	return &API{
		shopService: shopService,
	}
}

// RegisterRoutes регистрирует маршруты
func (api *API) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/api")
	{
		userGroup.POST("/sendCoin", api.SendCoin)
		userGroup.GET("/buy/:item", api.Buy)
		userGroup.GET("/info", api.Info)
	}
}
