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
		userGroup.GET("/user", api.Health)
		userGroup.GET("/buy/:item", api.middleware.ExtractUserID(), api.Buy)
	}
}

// Health проверка состояния
// @Summary Получить информацию о монетах
// @Description Получить информацию о монетах, инвентаре и истории транзакций
// @Tags shop
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200
// @Router /api/user [get]
func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
