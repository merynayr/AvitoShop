package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
)

// SendCoin обрабатывает HTTP-запрос на перевод монет
func (a *API) SendCoin(c *gin.Context) {
	var req struct {
		ToUser string `json:"toUser" binding:"required"`
		Amount int64  `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid JSON"})
		return
	}

	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	user, ok := userCtx.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user"})
		return
	}

	err := a.userService.SendCoins(c.Request.Context(), user.ID, req.ToUser, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
}
