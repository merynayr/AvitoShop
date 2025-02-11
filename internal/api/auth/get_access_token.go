package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var req struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GetAccessToken обрабатывает HTTP-запрос на получение access токена
func (a *API) GetAccessToken(c *gin.Context) {

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := a.authService.GetAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
