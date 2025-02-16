package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
)

// Login аутентификация и получение JWT-токена
// @Summary Аутентификация
// @Description Аутентифицирует пользователя и возвращает JWT-токен
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body model.AuthRequest true "Данные для аутентификации"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} sys.ErrorResponse
// @Failure 401 {object} sys.ErrorResponse
// @Failure 500 {object} sys.ErrorResponse
// @Router /api/auth [post]
func (a *API) Login(c *gin.Context) {
	var req model.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		sys.HandleError(c, sys.InvalidRequestError)
		return
	}

	authResponse, err := a.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		sys.HandleError(c, err)
		return
	}

	a.setCookies(c, authResponse.RefreshToken, authResponse.AccessToken)

	c.JSON(http.StatusOK, authResponse)
}
