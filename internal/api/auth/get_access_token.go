package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/sys"
)

// GetAccessToken обрабатывает HTTP-запрос на получение access токена
func (a *API) GetAccessToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		sys.HandleError(c, sys.InvalidRequestError)
		return
	}

	token, err := a.authService.GetAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		sys.HandleError(c, sys.InvalidAccessTokenError)
		return
	}

	a.setCookies(c, "", token)

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
