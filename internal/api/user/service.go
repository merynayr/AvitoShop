package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/middleware"
	"github.com/merynayr/AvitoShop/internal/service"
)

// API user структура
// объект сервисного слоя (его интерфейса)
type API struct {
	userService service.UserService
	middleware  middleware.UserMiddleware
}

// NewAPI возвращает новый объект имплементации API-слоя
func NewAPI(userService service.UserService, middleware middleware.UserMiddleware) *API {
	return &API{
		userService: userService,
		middleware:  middleware,
	}
}

// RegisterRoutes регистрирует маршруты
func (api *API) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/api")
	{
		userGroup.GET("/user", api.Health)
		userGroup.POST("/sendCoin", api.SendCoin)
		userGroup.GET("/buy/:item", api.Buy)
		userGroup.GET("/info", api.Info)
	}
}

// Health проверяет доступность сервиса
func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
