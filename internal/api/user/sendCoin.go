package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/sys/codes"
)

// SendCoin отправляет монеты другому пользователю
// @Summary Отправить монеты другому пользователю
// @Description Переводит указанное количество монет другому пользователю
// @Tags shop
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param body body model.SendCoinRequest true "Данные для перевода монет"
// @Success 200
// @Failure 400 {object} sys.ErrorResponse
// @Failure 401 {object} sys.ErrorResponse
// @Failure 500 {object} sys.ErrorResponse
// @Router /api/sendCoin [post]
func (a *API) SendCoin(c *gin.Context) {
	var req model.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sys.HandleError(c, sys.NewCommonError("invalid request", codes.BadRequest))
		return
	}

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

	err := a.userService.SendCoins(c.Request.Context(), user.ID, req)
	if err != nil {
		sys.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
}
