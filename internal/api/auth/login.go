package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

// Login обрабатывает HTTP-запрос на авторизацию
func (a *API) Login(c *gin.Context) {
	var req model.UserInfo

	if err := c.ShouldBindJSON(&req); err != nil {
		sys.HandleError(c, sys.NewCommonError("invalid request format", codes.BadRequest))
		return
	}

	refreshToken, err := a.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		fmt.Println(err)
		sys.HandleError(c, err)
		return
	}

	accessToken, err := a.authService.GetAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		sys.HandleError(c, err)
		return
	}

	a.setCookies(c, refreshToken, accessToken)

	c.JSON(http.StatusOK, gin.H{
		"refresh_token": refreshToken,
		"access_token":  accessToken,
	})
}
