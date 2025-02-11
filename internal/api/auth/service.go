package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/service"
)

// API auth структура
type API struct {
	authService service.AuthService
}

// NewAPI возвращает новый объект имплементации API-слоя auth
func NewAPI(authService service.AuthService) *API {
	return &API{
		authService: authService,
	}
}

// RegisterRoutes регистрирует маршруты
func (api *API) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/api")
	{
		authGroup.POST("/auth", api.Login)
		authGroup.POST("/access-token", api.GetAccessToken)
		authGroup.POST("/refresh-token", api.GetRefreshToken)
	}
}
