package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/sys"
)

// GetRefreshToken обрабатывает HTTP-запрос на получение refresh токена
func (a *API) GetRefreshToken(c *gin.Context) {
	oldRefreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		sys.HandleError(c, sys.InvalidRequestError)
		return
	}

	token, err := a.authService.GetRefreshToken(c.Request.Context(), oldRefreshToken)
	if err != nil {
		sys.HandleError(c, sys.InvalidRefreshTokenError)
		return
	}

	a.setCookies(c, token, "")
	c.JSON(http.StatusOK, gin.H{"refresh_token": token})
}
