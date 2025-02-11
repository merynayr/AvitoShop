package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRefreshToken обрабатывает HTTP-запрос на получение refresh токена
func (a *API) GetRefreshToken(c *gin.Context) {
	var req struct {
		OldRefreshToken string `json:"old_refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := a.authService.GetRefreshToken(c.Request.Context(), req.OldRefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"refresh_token": token})
}
