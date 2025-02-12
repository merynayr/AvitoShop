package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
)

// Info обрабатывает HTTP-запрос на получение информации
func (a *API) Info(c *gin.Context) {
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

	userInfo, err := a.userService.GetUserInfo(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
