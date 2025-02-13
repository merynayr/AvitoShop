package user

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/middleware"
	"github.com/merynayr/AvitoShop/internal/service"
)

// API user структура
// объект сервисного слоя (его интерфейса)
type API struct {
	userService service.UserService
	middleware  middleware.UserMiddlware
}

// NewAPI возвращает новый объект имплементации API-слоя
func NewAPI(userService service.UserService, middleware middleware.UserMiddlware) *API {
	return &API{
		userService: userService,
		middleware:  middleware,
	}
}

// RegisterRoutes регистрирует маршруты
func (api *API) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/api")
	{
		userGroup.POST("/sendCoin", api.middleware.ExtractUserID(), api.SendCoin)
		userGroup.GET("/buy/:item", api.middleware.ExtractUserID(), api.Buy)
		userGroup.GET("/info", api.middleware.ExtractUserID(), api.Info)
	}
}
