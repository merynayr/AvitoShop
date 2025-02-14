package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/sys"
)

// Info получает информацию о монетах, инвентаре и истории транзакций
// @Summary Получить информацию о монетах, инвентаре и истории транзакций
// @Description Возвращает баланс монет, содержимое инвентаря и историю транзакций пользователя
// @Tags shop
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} model.InfoResponse
// @Failure 400 {object} sys.ErrorResponse
// @Failure 401 {object} sys.ErrorResponse
// @Failure 500 {object} sys.ErrorResponse
// @Router /api/info [get]
func (a *API) Info(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		sys.HandleError(c, sys.UserNotFoundError)
		return
	}

	user, ok := userCtx.(*model.User)
	if !ok {
		sys.HandleError(c, fmt.Errorf(sys.ErrInvalidUser))
		return
	}

	userInfo, err := a.userService.GetUserInfo(c, user)
	if err != nil {
		sys.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
