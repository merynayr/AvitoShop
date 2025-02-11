package shop

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/service"
)

// API shop структура
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
		userGroup.GET("/shop", api.Health)
	}
}

// Health проверка состояния
func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
