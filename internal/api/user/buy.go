package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
)

// Buy обрабатывает HTTP-запрос на покупку мерча
func (a *API) Buy(c *gin.Context) {
	item := c.Param("item")
	item = strings.Trim(item, " \t\n\r")

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

	err := a.userService.Buy(c, user, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
