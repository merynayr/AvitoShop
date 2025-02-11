package converter

import (
	"github.com/merynayr/AvitoShop/internal/model"
	modelRepo "github.com/merynayr/AvitoShop/internal/repository/user/model"
)

// ToUserFromRepo конвертирует модель пользователя репо слоя в
// модель сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Coins:    user.Coins,
	}
}

// ToUserInfoFromRepo конвертирует модель пользователя репо слоя в
// модель сервисного слоя для авторизации
func ToUserInfoFromRepo(user *modelRepo.UserInfo) *model.UserInfo {
	if user == nil {
		return nil
	}

	return &model.UserInfo{
		Username: user.Username,
		Password: user.Password,
	}
}
