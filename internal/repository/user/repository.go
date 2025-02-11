package user

import (
	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/repository"
)

// Репозиторий
type repo struct {
	db db.Client
}

// NewRepository возвращает объект репозитория
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{
		db: db,
	}
}
