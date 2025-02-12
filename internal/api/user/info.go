package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/merynayr/AvitoShop/internal/model"
)

// Info обрабатывает HTTP-запрос на авторизацию
func (a *API) Info(c *gin.Context) {
	var req model.User
	_ = req

	c.JSON(http.StatusOK, gin.H{})
}
