package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

// Buy покупает предмет за монеты
// @Summary Купить предмет
// @Description Покупает указанный предмет за монеты
// @Tags shop
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param item path string true "Название предмета"
// @Success 200
// @Failure 400 {object} sys.ErrorResponse
// @Failure 401 {object} sys.ErrorResponse
// @Failure 500 {object} sys.ErrorResponse
// @Router /api/buy/{item} [get]
func (a *API) Buy(c *gin.Context) {
	item := c.Param("item")
	item = strings.Trim(item, " \t\n\r")

	userCtx, exists := c.Get("user")
	if !exists {
		sys.HandleError(c, sys.NewCommonError("user not found", codes.Unauthorized))
		return
	}

	user, ok := userCtx.(*model.User)
	if !ok {
		sys.HandleError(c, fmt.Errorf("invalid user"))
		return
	}

	err := a.userService.Buy(c, user, item)
	if err != nil {
		sys.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
