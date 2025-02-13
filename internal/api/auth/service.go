package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/service"
)

// API auth структура
type API struct {
	authService service.AuthService
	authConfig  config.AuthConfig
}

// NewAPI возвращает новый объект имплементации API-слоя auth
func NewAPI(authService service.AuthService, authConfig config.AuthConfig) *API {
	return &API{
		authService: authService,
		authConfig:  authConfig,
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

// setCookies устанавливают токены в куки
func (api *API) setCookies(c *gin.Context, refreshToken string, accessToken string) {
	if len(refreshToken) > 0 {
		c.SetCookie("refresh_token", refreshToken, int(api.authConfig.RefreshTokenExp()*2), "/api", "", false, true)
	}
	if len(accessToken) > 0 {
		c.SetCookie("access_token", accessToken, int(api.authConfig.AccessTokenExp()*2), "/api", "", false, true)
	}
}
