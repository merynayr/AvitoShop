package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

// GetRefreshToken обрабатывает HTTP-запрос на получение refresh токена
func (a *API) GetRefreshToken(c *gin.Context) {
	var req struct {
		OldRefreshToken string `json:"old_refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		sys.HandleError(c, sys.NewCommonError("invalid request", codes.BadRequest))
		return
	}

	token, err := a.authService.GetRefreshToken(c.Request.Context(), req.OldRefreshToken)
	if err != nil {
		sys.HandleError(c, sys.NewCommonError("invalid refresh token", codes.Unauthorized))
		return
	}

	c.JSON(http.StatusOK, gin.H{"refresh_token": token})
}
