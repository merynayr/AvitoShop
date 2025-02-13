package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

var req struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GetAccessToken обрабатывает HTTP-запрос на получение access токена
func (a *API) GetAccessToken(c *gin.Context) {

	if err := c.ShouldBindJSON(&req); err != nil {
		sys.HandleError(c, sys.NewCommonError("invalid request", codes.BadRequest))
		return
	}

	token, err := a.authService.GetAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		sys.HandleError(c, sys.NewCommonError("invalid access token", codes.Unauthorized))
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
