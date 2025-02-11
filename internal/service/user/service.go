package user

import (
	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/repository"
	"github.com/merynayr/AvitoShop/internal/service"
)

// Структура сервисного слоя с объектами репо слоя
// и транзакционного менеджера
type userService struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService возвращает объект сервисного слоя
func NewService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &userService{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
